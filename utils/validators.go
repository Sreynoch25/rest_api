package utils

import (
	"errors"
	"fiber-crud/models"
	"regexp"
)

func ValidateEmail(email string) bool {
	pattern :=  `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	
	return re.MatchString(email)
}

func ValidatePhone(phone string) bool {
	pattern :=  `^\+?[0-9]{10,15}$`
	re := regexp.MustCompile(pattern)
	
	return re.MatchString(phone)
}

func ValidateUser(user *models.User) error {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
	}

	if !ValidateEmail(user.Email) {
        return errors.New("invalid email format")
    }

	if len(user.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
	
    // if user.Phone != "" && !ValidatePhone(user.Phone) {
    //     return errors.New("invalid phone format")
    // }

	return nil
}