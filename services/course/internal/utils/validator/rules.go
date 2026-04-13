package validator

import "github.com/ritchieridanko/klasshub/services/course/internal/constants"

var (
	roleAllowedSubdomains = map[string]string{
		constants.RoleAdministrator: constants.SubdomainAdmin,
		constants.RoleSchool:        constants.SubdomainAdmin,
		constants.RoleInstructor:    constants.SubdomainLMS,
		constants.RoleStudent:       constants.SubdomainLMS,
	}
)

const (
	courseDescMaxLength     int = 1000
	courseNameMaxLength     int = 100
	courseNameMinLength     int = 2
	schoolCourseIDMaxLength int = 100
)
