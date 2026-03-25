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
			constants.SubdomainLMS:   {},
			constants.SubdomainAdmin: {},
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
}
