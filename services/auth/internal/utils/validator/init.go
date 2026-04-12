package validator

import "fmt"

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Email(value string) (bool, string) {
	if !rgxEmail.MatchString(value) {
		return false, fmt.Sprintf("Email is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) Identifier(value string) (bool, string) {
	isEmail, _ := v.Email(value)
	isUsername, _ := v.Username(value)
	if !isEmail && !isUsername {
		return false, fmt.Sprintf("Email/Username is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) Password(value string) (bool, string) {
	length := len(value)
	if length < passwordMinLength {
		return false, fmt.Sprintf("Password must be at least %d characters", passwordMinLength)
	}
	if length > passwordMaxLength {
		return false, fmt.Sprintf("Password must not exceed %d characters", passwordMaxLength)
	}
	if !rgxLowercase.MatchString(value) {
		return false, "Password must include at least one lowercase letter"
	}
	if !rgxUppercase.MatchString(value) {
		return false, "Password must include at least one uppercase letter"
	}
	if !rgxNumber.MatchString(value) {
		return false, "Password must include at least one number"
	}
	if !rgxSpecialChars.MatchString(value) {
		return false, fmt.Sprintf(
			"Password must include at least one special character: %s",
			specialChars,
		)
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

func (v *Validator) Username(value string) (bool, string) {
	if !rgxUsername.MatchString(value) {
		return false, fmt.Sprintf("Username is invalid: %s", value)
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
