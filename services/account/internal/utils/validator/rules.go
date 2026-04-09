package validator

import "github.com/ritchieridanko/klasshub/services/account/internal/constants"

var (
	roleAllowedSubdomains = map[string]string{
		constants.RoleAdministrator: constants.SubdomainAdmin,
		constants.RoleSchool:        constants.SubdomainAdmin,
		constants.RoleInstructor:    constants.SubdomainLMS,
		constants.RoleStudent:       constants.SubdomainLMS,
	}
)
