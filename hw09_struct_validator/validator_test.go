package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstraints(t *testing.T) {
	for _, tt := range []struct {
		raw       string
		expected  map[string]string
		expectErr bool
	}{
		{raw: `min:10`, expected: map[string]string{"min": "10"}},
		{raw: `min:10|max:20`, expected: map[string]string{"min": "10", "max": "20"}},
		{raw: `nested`, expected: map[string]string{"nested": ""}},
		{raw: `min:10|wrong-format`, expected: nil, expectErr: true},
		{raw: "", expected: nil, expectErr: true},
	} {
		result, err := constraints(tt.raw)
		if tt.expectErr {
			require.NotNil(t, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tt.expected, result)
	}
}

type UserRole string

// Test the function on different structures and other types.
type (
	Group struct {
		ID   string `validate:"len:36"`
		Name string
	}

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
		Group  Group `validate:"nested"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "8bc20e39-1d1e-46f0-b762-b90e67746b3c",
				Age:    18,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"+1234567890", "+2345678901"},
				Group: Group{
					ID:   "7aeddec6-ce81-4072-a4ab-63913551a8ad",
					Name: "group name",
				},
			},
			expectedErr: nil,
		},
		{
			in: User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLenNotEqual{Len: 0, Expected: 36}},
				ValidationError{Field: "Age", Err: ErrLessThan{Value: 0, Min: 18}},
				ValidationError{
					Field: "Email",
					Err:   ErrRegexpNotMatch{Value: "", Regexp: regexp.MustCompile(`^\w+@\w+\.\w+$`)},
				},
				ValidationError{Field: "Role", Err: ErrStrNotIn{Value: "", Items: []string{"admin", "stuff"}}},
				ValidationError{
					Field: "Group",
					Err:   ValidationErrors{ValidationError{Field: "ID", Err: ErrLenNotEqual{Len: 0, Expected: 36}}},
				},
			},
		},
		{
			in: User{
				ID:     "123456789-123456789-123456789-123456789-",
				Age:    60,
				Email:  "not-an-email",
				Role:   "wrong-role",
				Phones: []string{"123", "1234"},
				Group: Group{
					ID: "1234567",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLenNotEqual{Len: 40, Expected: 36}},
				ValidationError{Field: "Age", Err: ErrGreaterThan{Value: 60, Max: 50}},
				ValidationError{
					Field: "Email",
					Err:   ErrRegexpNotMatch{Value: "not-an-email", Regexp: regexp.MustCompile(`^\w+@\w+\.\w+$`)},
				},
				ValidationError{
					Field: "Role",
					Err:   ErrStrNotIn{Value: "wrong-role", Items: []string{"admin", "stuff"}},
				},
				ValidationError{Field: "Phones", Err: ErrLenNotEqual{Len: 3, Expected: 11}},
				ValidationError{Field: "Phones", Err: ErrLenNotEqual{Len: 4, Expected: 11}},
				ValidationError{
					Field: "Group",
					Err:   ValidationErrors{ValidationError{Field: "ID", Err: ErrLenNotEqual{Len: 7, Expected: 36}}},
				},
			},
		},

		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in: App{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLenNotEqual{Len: 0, Expected: 5}},
			},
		},
		{
			in: App{Version: "123"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLenNotEqual{Len: 3, Expected: 5}},
			},
		},
		{
			in: App{Version: "1234567"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLenNotEqual{Len: 7, Expected: 5}},
			},
		},

		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("any"),
				Payload:   []byte("any"),
				Signature: []byte("any"),
			},
			expectedErr: nil,
		},

		{
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			in: Response{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntNotIn{Value: 0, Items: []int{200, 404, 500}}},
			},
		},
		{
			in: Response{Code: 201},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntNotIn{Value: 201, Items: []int{200, 404, 500}}},
			},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.Equal(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
