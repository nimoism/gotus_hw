package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	tagName     = "validate"
	tagSep      = "|"
	tagValueSep = ":"

	tagMin    = "min"
	tagMax    = "max"
	tagIn     = "in"
	tagInSep  = ","
	tagLen    = "len"
	tagRegexp = "regexp"
	tagNested = "nested"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	var vErrs ValidationErrors
	if ok := errors.As(v.Err, &vErrs); ok {
		return fmt.Sprintf("%v: { %v }", v.Field, vErrs)
	}
	return fmt.Sprintf("%v: %v", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStrings := make([]string, 0, len(v))
	for _, fieldErr := range v {
		errStrings = append(errStrings, fieldErr.Error())
	}
	return strings.Join(errStrings, ", ")
}

func Validate(v interface{}) error {
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	errs := ValidationErrors{}
	for i := 0; i < vType.NumField(); i++ {
		fStruct := vType.Field(i)
		fValue := vValue.Field(i)
		tag, ok := fStruct.Tag.Lookup(tagName)
		if !ok {
			continue
		}
		constraintsMap, err := constraints(tag)
		if err != nil {
			return err
		}
		for cName, constraint := range constraintsMap {
			if err := validateValue(fStruct.Name, fStruct.Type, fValue, cName, constraint); err != nil {
				var ok bool
				var vErr ValidationError
				if ok = errors.As(err, &vErr); ok {
					errs = append(errs, vErr)
					continue
				}
				var nestedErrs ValidationErrors
				if ok = errors.As(err, &nestedErrs); ok {
					errs = append(errs, nestedErrs...)
					continue
				}
				return err
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateValue(fName string, fType reflect.Type, fValue reflect.Value, constraintName, constraint string) error {
	switch fType.Kind() { // nolint: exhaustive
	case reflect.Int:
		validate, ok := intValidates[constraintName]
		if !ok {
			return fmt.Errorf("unknown tag '%v'", constraintName)
		}
		return validate(fName, constraint, int(fValue.Int()))
	case reflect.String:
		validate, ok := strValidates[constraintName]
		if !ok {
			return fmt.Errorf("unknown tag '%v'", constraintName)
		}
		return validate(fName, constraint, fValue.String())
	case reflect.Slice:
		nestedErrs := ValidationErrors{}
		for i := 0; i < fValue.Len(); i++ {
			if err := validateValue(fName, fType.Elem(), fValue.Index(i), constraintName, constraint); err != nil {
				var vErr ValidationError
				if errors.As(err, &vErr) {
					nestedErrs = append(nestedErrs, vErr)
					continue
				}
				return err
			}
		}
		return nestedErrs
	case reflect.Struct:
		if constraintName != tagNested {
			return fmt.Errorf("unknown tag '%v'", constraintName)
		}
		var nestedErrs ValidationErrors
		if err := Validate(fValue.Interface()); err != nil {
			if errors.As(err, &nestedErrs) {
				return ValidationError{Field: fName, Err: err}
			}
			return err
		}
	}
	return nil
}

func constraints(tag string) (map[string]string, error) {
	rawConstraints := strings.Split(tag, tagSep)
	constraintsMap := make(map[string]string, len(rawConstraints))
	for _, rawConstraint := range rawConstraints {
		rawParts := strings.SplitN(rawConstraint, tagValueSep, 2)
		if len(rawParts) < 1 {
			return nil, fmt.Errorf("empty validate min")
		}
		name := rawParts[0]
		if name == tagNested {
			constraintsMap[name] = ""
			continue
		}
		if len(rawParts) < 2 {
			return nil, fmt.Errorf("wrong validate min format '%v'", rawConstraint)
		}
		constraintsMap[name] = rawParts[1]
	}
	return constraintsMap, nil
}
