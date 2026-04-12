package policies

import "github.com/ritchieridanko/klasshub/services/auth/internal/constants"

type AuthPolicy struct {
	requireAuth         bool
	requireSchool       bool
	requireVerification bool
	roles               map[string]struct{}
	subdomains          map[string]struct{}
}

func (p *AuthPolicy) RequireAuth() bool {
	return p.requireAuth
}

func (p *AuthPolicy) RequireSchool() bool {
	return p.requireSchool
}

func (p *AuthPolicy) RequireVerification() bool {
	return p.requireVerification
}

func (p *AuthPolicy) RequireRole() bool {
	return len(p.roles) > 0
}

func (p *AuthPolicy) RequireSubdomain() bool {
	return len(p.subdomains) > 0
}

func (p *AuthPolicy) IsRoleAuthorized(role string) bool {
	_, exists := p.roles[role]
	return exists
}

func (p *AuthPolicy) IsSubdomainAuthorized(subdomain string) bool {
	_, exists := p.subdomains[subdomain]
	return exists
}

var AuthPolicies map[string]AuthPolicy = map[string]AuthPolicy{
	"/auth.v1.AuthService/Login": AuthPolicy{
		requireAuth:         false,
		requireSchool:       false,
		requireVerification: false,
		roles:               map[string]struct{}{},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
			constants.SubdomainLMS:   {},
		},
	},
	"/auth.v1.AuthService/Logout": AuthPolicy{
		requireAuth:         true,
		requireSchool:       false,
		requireVerification: false,
		roles:               map[string]struct{}{},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
			constants.SubdomainLMS:   {},
		},
	},
	"/auth.v1.AuthService/CreateSchoolAuth": AuthPolicy{
		requireAuth:         false,
		requireSchool:       false,
		requireVerification: false,
		roles:               map[string]struct{}{},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/CreateUserAuth": AuthPolicy{
		requireAuth:         true,
		requireSchool:       true,
		requireVerification: true,
		roles: map[string]struct{}{
			constants.RoleAdministrator: {},
			constants.RoleSchool:        {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/UpdateSchool": AuthPolicy{
		requireAuth:         true,
		requireSchool:       true,
		requireVerification: true,
		roles: map[string]struct{}{
			constants.RoleSchool: {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/ChangePassword": AuthPolicy{
		requireAuth:         true,
		requireSchool:       true,
		requireVerification: true,
		roles: map[string]struct{}{
			constants.RoleAdministrator: {},
			constants.RoleInstructor:    {},
			constants.RoleSchool:        {},
			constants.RoleStudent:       {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
			constants.SubdomainLMS:   {},
		},
	},
	"/auth.v1.AuthService/ResendVerification": AuthPolicy{
		requireAuth:         true,
		requireSchool:       true,
		requireVerification: false,
		roles: map[string]struct{}{
			constants.RoleSchool: {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/VerifyEmail": AuthPolicy{
		requireAuth:         true,
		requireSchool:       true,
		requireVerification: false,
		roles: map[string]struct{}{
			constants.RoleSchool: {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/RotateAuthToken": AuthPolicy{
		requireAuth:         false,
		requireSchool:       false,
		requireVerification: false,
		roles:               map[string]struct{}{},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
			constants.SubdomainLMS:   {},
		},
	},
	"/auth.v1.AuthService/IsEmailAvailable": AuthPolicy{
		requireAuth:         false,
		requireSchool:       false,
		requireVerification: false,
		roles:               map[string]struct{}{},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
	"/auth.v1.AuthService/IsUsernameAvailable": AuthPolicy{
		requireAuth:         false,
		requireSchool:       false,
		requireVerification: false,
		roles: map[string]struct{}{
			constants.RoleAdministrator: {},
			constants.RoleSchool:        {},
		},
		subdomains: map[string]struct{}{
			constants.SubdomainAdmin: {},
		},
	},
}
