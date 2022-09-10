package hw10programoptimization

import (
	"bufio"
	"fmt"
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
	emailSepo      = "@"
	emailFieldName = "Email"
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	emails, err := getEmails(r)
	if err != nil {
		return nil, fmt.Errorf("getting email error: %w", err)
	}
	return countDomains(emails, domain)
}

func getEmails(r io.Reader) ([]string, error) {
	emails := make([]string, 0, 128000)
	scanner := bufio.NewScanner(r)

	var parser fastjson.Parser

	for scanner.Scan() {
		userJSON, err := parser.Parse(scanner.Text())
		if err != nil {
			return nil, err
		}
		emails = append(emails, string(userJSON.GetStringBytes(emailFieldName)))
	}
	return emails, nil
}

func countDomains(emails []string, domain string) (DomainStat, error) {
	count := make(DomainStat)
	domain = "." + domain
	for _, email := range emails {
		matched := strings.Contains(email, domain)

		if matched {
			count[strings.ToLower(strings.SplitN(email, emailSepo, 2)[1])]++
		}
	}
	return count, nil
}
