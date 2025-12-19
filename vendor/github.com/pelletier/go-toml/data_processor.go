
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if p.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if len(p.Username) < 3 || len(p.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func NormalizeProfile(p UserProfile) UserProfile {
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	p.Username = strings.TrimSpace(p.Username)
	return p
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	p = NormalizeProfile(p)
	if err := ValidateProfile(p); err != nil {
		return UserProfile{}, err
	}
	return p, nil
}