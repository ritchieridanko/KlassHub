package validator

import "github.com/ritchieridanko/klasshub/services/user/internal/constants"

var (
	roles = map[string]struct{}{
		constants.RoleAdministrator: {},
		constants.RoleInstructor:    {},
		constants.RoleStudent:       {},
	}

	roleAllowedSubdomains = map[string]string{
		constants.RoleAdministrator: constants.SubdomainAdmin,
		constants.RoleSchool:        constants.SubdomainAdmin,
		constants.RoleInstructor:    constants.SubdomainLMS,
		constants.RoleStudent:       constants.SubdomainLMS,
	}
)

const (
	birthplaceMaxLength   int = 100
	birthplaceMinLength   int = 2
	nameMaxLength         int = 100
	nameMinLength         int = 3
	schoolUserIDMaxLength int = 50
)
