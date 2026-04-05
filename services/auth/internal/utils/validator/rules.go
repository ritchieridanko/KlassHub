package validator

import (
	"regexp"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
)

var (
	rgxEmail        *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase    *regexp.Regexp = regexp.MustCompile(`[a-z]`)
	rgxNumber       *regexp.Regexp = regexp.MustCompile(`[0-9]`)
	rgxSpecialChars *regexp.Regexp = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	rgxUppercase    *regexp.Regexp = regexp.MustCompile(`[A-Z]`)
	rgxUsername     *regexp.Regexp = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9]|(?:[._][a-z0-9])){7,24}$`)

	roleAllowedSubdomains = map[string]string{
		constants.RoleAdministrator: constants.SubdomainAdmin,
		constants.RoleSchool:        constants.SubdomainAdmin,
		constants.RoleInstructor:    constants.SubdomainLMS,
		constants.RoleStudent:       constants.SubdomainLMS,
	}
)

const (
	passwordMaxLength int = 50
	passwordMinLength int = 8

	specialChars string = `!@#$%^&*(),.?":{}|<>`
)
