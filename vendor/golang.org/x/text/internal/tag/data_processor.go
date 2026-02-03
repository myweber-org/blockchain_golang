package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Username  string
	Email     string
	Age       int
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(profile UserProfile) error {
	if profile.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(profile.Username) < 3 || len(profile.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}
	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	transformed.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	return transformed
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}
	transformed := TransformProfile(profile)
	return transformed, nil
}