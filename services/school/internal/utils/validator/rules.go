package validator

import (
	"regexp"

	"github.com/ritchieridanko/klasshub/services/school/internal/constants"
)

var (
	rgxEmail    *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxNPSN     *regexp.Regexp = regexp.MustCompile(`^\d{8}$`)
	rgxPhone    *regexp.Regexp = regexp.MustCompile(`^0\d{7,15}$`)
	rgxPostcode *regexp.Regexp = regexp.MustCompile(`^\d{5}$`)

	roleAllowedSubdomains = map[string]string{
		constants.RoleAdministrator: constants.SubdomainAdmin,
		constants.RoleSchool:        constants.SubdomainAdmin,
		constants.RoleInstructor:    constants.SubdomainLMS,
		constants.RoleStudent:       constants.SubdomainLMS,
	}
)

const (
	schoolNameMaxLength int = 100
	schoolNameMinLength int = 2
	streetMaxLength     int = 100
	streetMinLength     int = 2
)
