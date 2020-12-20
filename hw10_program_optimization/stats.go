package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var ErrWrongEmail = errors.New("wrong email")
var ErrWrongDomain = errors.New("wrong domain")

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var topDomainRegex = regexp.MustCompile(`^[[:alpha:]]{2,}$`)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if !topDomainRegex.MatchString(domain) {
		return nil, ErrWrongDomain
	}
	domain = `.` + strings.ToLower(domain)
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	ji := jsoniter.ConfigFastest
	var user User
	for i := 0; scanner.Scan(); i++ {
		user = User{}
		if err := ji.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		email := strings.ToLower(user.Email)
		if !strings.HasSuffix(email, domain) {
			continue
		}
		split := strings.Split(email, "@")
		if len(split) != 2 {
			return nil, ErrWrongEmail
		}
		result[split[1]]++
	}
	return result, nil
}
