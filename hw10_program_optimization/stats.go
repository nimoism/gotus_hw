package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var ErrWrongEmail = errors.New("wrong email")
var ErrEmptyDomain = errors.New("empty domain")

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type users [100_000]User
type DomainStat map[string]int

// nolint: errorlint
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	ji := jsoniter.ConfigFastest
	var u User
	for i := 0; scanner.Scan(); i++ {
		if err = ji.Unmarshal(scanner.Bytes(), &u); err != nil {
			return
		}
		result[i] = u
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}
	result := make(DomainStat)
	dottedDomain := `.` + domain
	for _, user := range u {
		if !strings.HasSuffix(user.Email, dottedDomain) {
			continue
		}
		split := strings.Split(user.Email, "@")
		if len(split) != 2 {
			return nil, ErrWrongEmail
		}
		result[strings.ToLower(split[1])]++
	}
	return result, nil
}
