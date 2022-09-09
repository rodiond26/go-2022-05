package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

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

var (
	ErrEmptyDomain = errors.New("domain is empty")
	emailSepo      = "@"
	emailFieldName = "Email"
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}
	return countDomains(r, domain)
}

type users [100_000]User

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	rr := bufio.NewReader(r)

	for {
		userJson, _, err := rr.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		email := fastjson.GetString(userJson, emailFieldName)
		if strings.Contains(email, domain) {
			result[strings.ToLower(strings.SplitN(email, emailSepo, 2)[1])]++
		}
	}
	return result, nil
}
