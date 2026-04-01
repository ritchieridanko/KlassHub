package models

type CreateSchoolAuthReq struct {
	Email    string
	Password string
}

type LoginReq struct {
	Identifier string
	Password   string
}

type VerifyEmailReq struct {
	VerificationToken string
	RefreshToken      string
}
