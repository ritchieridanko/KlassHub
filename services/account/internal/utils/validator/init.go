package validator

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) RoleAllowedSubdomain(role, subdomain string) bool {
	sd, ok := roleAllowedSubdomains[role]
	if !ok || subdomain != sd {
		return false
	}
	return true
}
