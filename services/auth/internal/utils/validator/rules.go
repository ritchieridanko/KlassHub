package validator

import "regexp"

var (
	rgxEmail        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase    = regexp.MustCompile(`[a-z]`)
	rgxNumber       = regexp.MustCompile(`[0-9]`)
	rgxSpecialChars = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	rgxUppercase    = regexp.MustCompile(`[A-Z]`)
	rgxUsername     = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9]|(?:[._][a-z0-9])){7,24}$`)
)

const (
	passwordMaxLength  int = 50
	passwordMinLength  int = 8
	userAgentMaxLength int = 512

	specialChars string = `!@#$%^&*(),.?":{}|<>`
)
