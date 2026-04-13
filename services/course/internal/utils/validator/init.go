package validator

import (
	"fmt"
	"net/url"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) CourseDesc(value string) (bool, string) {
	if len(value) > courseDescMaxLength {
		return false, fmt.Sprintf("Course Description must not exceed %d characters", courseDescMaxLength)
	}
	return true, ""
}

func (v *Validator) CourseName(value string) (bool, string) {
	length := len(value)
	if length < courseNameMinLength {
		return false, fmt.Sprintf("Course Name must be at least %d characters", courseNameMinLength)
	}
	if length > courseNameMaxLength {
		return false, fmt.Sprintf("Course Name must not exceed %d characters", courseNameMaxLength)
	}
	return true, ""
}

func (v *Validator) SchoolCourseID(value string) (bool, string) {
	if len(value) > schoolCourseIDMaxLength {
		return false, fmt.Sprintf("School Course ID must not exceed %d characters", schoolCourseIDMaxLength)
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
