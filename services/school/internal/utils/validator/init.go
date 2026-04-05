package validator

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ritchieridanko/klasshub/shared/data"
)

type Validator struct {
	school struct {
		accreditations map[string]string
		levels         map[string]string
		ownerships     map[string]string
	}
}

func Init(sd *data.School) *Validator {
	return &Validator{
		school: struct {
			accreditations map[string]string
			levels         map[string]string
			ownerships     map[string]string
		}{
			accreditations: sd.Accreditations(),
			levels:         sd.Levels(),
			ownerships:     sd.Ownerships(),
		},
	}
}

func (v *Validator) Email(value string) (bool, string) {
	if !rgxEmail.MatchString(value) {
		return false, fmt.Sprintf("Email is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) NPSN(value string) (bool, string) {
	if !rgxNPSN.MatchString(value) {
		return false, fmt.Sprintf("NPSN is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) Phone(value string) (bool, string) {
	if !rgxPhone.MatchString(value) {
		return false, fmt.Sprintf("Phone is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) Postcode(value string) (bool, string) {
	if !rgxPostcode.MatchString(value) {
		return false, fmt.Sprintf("Postcode is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) SchoolAccreditation(value string) (bool, string) {
	_, ok := v.school.accreditations[value]
	if !ok {
		return false, fmt.Sprintf("School Accreditation is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) SchoolEstablishedAt(value time.Time) (bool, string) {
	if value.After(time.Now().UTC()) {
		return false, fmt.Sprintf("School Establishment Date is invalid: %s", value.Format("2 Jan 2006"))
	}
	return true, ""
}

func (v *Validator) SchoolLevel(value string) (bool, string) {
	_, ok := v.school.levels[value]
	if !ok {
		return false, fmt.Sprintf("School Level is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) SchoolName(value string) (bool, string) {
	length := len(value)
	if length < schoolNameMinLength {
		return false, fmt.Sprintf("School Name must be at least %d characters", schoolNameMinLength)
	}
	if length > schoolNameMaxLength {
		return false, fmt.Sprintf("School Name must not exceed %d characters", schoolNameMaxLength)
	}
	return true, ""
}

func (v *Validator) SchoolOwnership(value string) (bool, string) {
	_, ok := v.school.ownerships[value]
	if !ok {
		return false, fmt.Sprintf("School Ownership is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) Street(value string) (bool, string) {
	length := len(value)
	if length < streetMinLength {
		return false, fmt.Sprintf("Street must be at least %d characters", streetMinLength)
	}
	if length > streetMaxLength {
		return false, fmt.Sprintf("Street must not exceed %d characters", streetMaxLength)
	}
	return true, ""
}

func (v *Validator) URL(value string) (bool, string) {
	val, err := url.ParseRequestURI(value)
	if err != nil {
		return false, fmt.Sprintf("URL is invalid: %s", value)
	}
	if val.Scheme != "http" && val.Scheme != "https" {
		return false, fmt.Sprintf("URL is invalid: %s", value)
	}
	if val.Host == "" {
		return false, fmt.Sprintf("URL is invalid: %s", value)
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
