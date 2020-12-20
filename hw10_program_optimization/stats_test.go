// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty users", func(t *testing.T) {
		result, err := GetDomainStat(strings.NewReader(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty domain error", func(t *testing.T) {
		_, err := GetDomainStat(strings.NewReader(data), "")
		require.True(t, errors.Is(err, ErrWrongDomain))
	})

	t.Run("wrong domain error", func(t *testing.T) {
		for _, domain := range []string{"", "c", "c0m", "c_m"} {
			_, err := GetDomainStat(strings.NewReader(data), domain)
			require.True(t, errors.Is(err, ErrWrongDomain), "non wrong domain: %s", domain)
		}
	})

	t.Run("email only ends with domain", func(t *testing.T) {
		data := `{"Id":1,"Email":"aliquid_qui_ea@Browsedrive.com.gov"}`
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("wrong email without @", func(t *testing.T) {
		data := `{"Id":2,"Email":"broWsecat.com"}`
		_, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.True(t, errors.Is(err, ErrWrongEmail))
	})

	t.Run("wrong email too much @", func(t *testing.T) {
		data := `{"Id":1,"Email":"wr@ng@email.com"}`
		_, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.True(t, errors.Is(err, ErrWrongEmail))
	})

	t.Run("domain case insensitive", func(t *testing.T) {
		data := `{"Id":2,"Email":"mLynch@broWsecat.com"}`
		result, err := GetDomainStat(bytes.NewBufferString(data), "COM")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})

	t.Run("email domain case insensitive", func(t *testing.T) {
		data := `{"Id":2,"Email":"mLynch@broWsecat.COM"}`
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})

}
