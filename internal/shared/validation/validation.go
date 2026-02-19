package validation

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

func NonEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

func ValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}
