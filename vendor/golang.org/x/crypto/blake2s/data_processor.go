
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate string
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	now := time.Now()
	age := now.Year() - parsedDate.Year()

	if now.YearDay() < parsedDate.YearDay() {
		age--
	}

	if age < 0 {
		return 0, errors.New("birth date cannot be in the future")
	}

	return age, nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = NormalizeUsername(profile.Username)

	age, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	if age < 13 {
		return profile, errors.New("user must be at least 13 years old")
	}

	return profile, nil
}