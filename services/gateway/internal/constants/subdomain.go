package constants

const (
	SubdomainAdmin string = "admin"
	SubdomainLMS   string = "lms"
)

var (
	// Subdomains: Admin and LMS
	AllSubdomains []string = []string{
		SubdomainAdmin,
		SubdomainLMS,
	}
)
