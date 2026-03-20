package validator

import "regexp"

var (
	rgxEmail    *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxNPSN     *regexp.Regexp = regexp.MustCompile(`^\d{8}$`)
	rgxPhone    *regexp.Regexp = regexp.MustCompile(`^0\d{7,15}$`)
	rgxPostcode *regexp.Regexp = regexp.MustCompile(`^\d{5}$`)
)

const (
	schoolNameMaxLength int = 100
	schoolNameMinLength int = 2
	streetMaxLength     int = 100
	streetMinLength     int = 2
)
