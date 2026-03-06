package validator

import (
	"fmt"
	"net"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) IPAddress(value string) (bool, string) {
	if ip := net.ParseIP(value); ip == nil {
		return false, fmt.Sprintf("IP Address is invalid: %s", value)
	}
	return true, ""
}

func (v *Validator) UserAgent(value string) (bool, string) {
	if len(value) > userAgentMaxLength {
		return false, fmt.Sprintf("User Agent must not exceed %d characters", userAgentMaxLength)
	}
	return true, ""
}
