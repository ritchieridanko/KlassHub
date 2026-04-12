package models

type VerificationEmailMsg struct {
	Recipient         string
	VerificationToken string
	Role              string
}

type WelcomeEmailMsg struct {
	Recipient         string
	VerificationToken string
	Role              string
}
