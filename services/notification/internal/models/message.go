package models

type VerificationEmailMsg struct {
	Recipient         string
	VerificationToken string
}

type WelcomeEmailMsg struct {
	Recipient         string
	VerificationToken string
}
