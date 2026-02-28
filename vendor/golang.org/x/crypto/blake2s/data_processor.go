
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if p.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}

	if len(strings.TrimSpace(p.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformProfile(p *UserProfile) {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
}

func ProcessUserData(profiles []UserProfile) ([]UserProfile, error) {
	var validProfiles []UserProfile

	for _, profile := range profiles {
		if err := ValidateProfile(profile); err != nil {
			continue
		}

		TransformProfile(&profile)
		validProfiles = append(validProfiles, profile)
	}

	if len(validProfiles) == 0 {
		return nil, errors.New("no valid profiles found")
	}

	return validProfiles, nil
}