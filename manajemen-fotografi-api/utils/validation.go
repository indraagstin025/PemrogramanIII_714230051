package utils

import "regexp"

func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{8,15}$`)
	return re.MatchString(phone)
}


