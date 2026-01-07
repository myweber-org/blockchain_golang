
package data_processor

import (
	"errors"
	"strings"
	"time"
)

type UserData struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	Active    bool
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	
	if data.CreatedAt.After(time.Now()) {
		return errors.New("creation date cannot be in the future")
	}
	
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func GenerateUserReport(users []UserData) map[string]int {
	report := make(map[string]int)
	
	for _, user := range users {
		if user.Active {
			report["active"]++
		} else {
			report["inactive"]++
		}
		
		if strings.HasSuffix(user.Email, ".com") {
			report["dotcom_domains"]++
		}
	}
	
	report["total"] = len(users)
	return report
}

func FilterActiveUsers(users []UserData) []UserData {
	var activeUsers []UserData
	
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	
	return activeUsers
}

func CalculateAverageID(users []UserData) float64 {
	if len(users) == 0 {
		return 0
	}
	
	var sum int
	for _, user := range users {
		sum += user.ID
	}
	
	return float64(sum) / float64(len(users))
}