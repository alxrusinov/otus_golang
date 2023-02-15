package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//easyjson:json
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	count := 0
	for scanner.Scan() {
		var user User
		line := scanner.Bytes()
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		result[count] = user
		count++
	}

	if err = scanner.Err(); err != nil {
		return
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	reg, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		if matched := reg.Match([]byte(user.Email)); matched {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[key]
			num++
			result[key] = num
		}
	}
	return result, nil
}
