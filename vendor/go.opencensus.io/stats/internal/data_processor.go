
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.TrimSpace(transformed.Email)

	if transformed.Age < 18 {
		transformed.Active = false
	}

	return transformed
}

func ProcessUserData(input []byte) ([]byte, error) {
	var profile UserProfile
	if err := json.Unmarshal(input, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	transformed := TransformProfile(profile)
	output, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return output, nil
}

func main() {
	inputData := []byte(`{
		"id": 123,
		"username": "JohnDoe",
		"email": "JOHN@EXAMPLE.COM",
		"age": 25,
		"active": true
	}`)

	result, err := ProcessUserData(inputData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Processed user profile:")
	fmt.Println(string(result))
}
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    var sum float64

    for i := 0; i < windowSize; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(windowSize)

    for i := windowSize; i < len(data); i++ {
        sum = sum - data[i-windowSize] + data[i]
        result[i-windowSize+1] = sum / float64(windowSize)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := movingAverage(sampleData, window)
    fmt.Printf("Moving averages with window %d: %v\n", window, averages)
}
package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}