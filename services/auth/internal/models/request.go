package models

type CreateSchoolAuthReq struct {
	Email    string
	Password string
}

type CreateSessionReq struct {
	AuthID     int64
	SchoolID   int64
	Role       string
	IsVerified bool
}

type LoginReq struct {
	Identifier string
	Password   string
}

type RefreshSessionReq struct {
	AuthID       int64
	SchoolID     int64
	Role         string
	IsVerified   bool
	RefreshToken string
}

type VerifyEmailReq struct {
	VerificationToken string
	RefreshToken      string
}
