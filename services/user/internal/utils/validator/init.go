package validator

import (
	"fmt"
	"time"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Birthdate(value time.Time) (bool, string) {
	if value.After(time.Now().UTC()) {
		return false, fmt.Sprintf("Birthdate is invalid: %s", value.Format("2 Jan 2006"))
	}
	return true, ""
}

func (v *Validator) Birthplace(value string) (bool, string) {
	length := len(value)
	if length < birthplaceMinLength {
		return false, fmt.Sprintf("Birthplace must be at least %d characters", birthplaceMinLength)
	}
	if length > birthplaceMaxLength {
		return false, fmt.Sprintf("Birthplace must not exceed %d characters", birthplaceMaxLength)
	}
	return true, ""
}

func (v *Validator) Name(value string) (bool, string) {
	length := len(value)
	if length < nameMinLength {
		return false, fmt.Sprintf("Name must be at least %d characters", nameMinLength)
	}
	if length > nameMaxLength {
		return false, fmt.Sprintf("Name must not exceed %d characters", nameMaxLength)
	}
	return true, ""
}

func (v *Validator) Role(value string) (bool, string) {
	_, exists := roles[value]
	if !exists {
		return false, fmt.Sprintf("Role is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) SchoolUserID(value string) (bool, string) {
	if len(value) > schoolUserIDMaxLength {
		return false, fmt.Sprintf("School User ID must not exceed %d characters", schoolUserIDMaxLength)
	}
	return true, ""
}

func (v *Validator) Sex(value string) (bool, string) {
	if value != "male" && value != "female" {
		return false, fmt.Sprintf("Sex is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) RoleAllowedSubdomain(role, subdomain string) bool {
	sd, ok := roleAllowedSubdomains[role]
	if !ok || subdomain != sd {
		return false
	}
	return true
}
