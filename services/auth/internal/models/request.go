package models

type ChangePasswordReq struct {
	OldPassword string
	NewPassword string
}

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

type CreateUserAuthReq struct {
	Email    *string
	Username *string
	Password string
	Role     string
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

type UpdateSchoolReq struct {
	SchoolID     int64
	RefreshToken string
}

type VerifyEmailReq struct {
	VerificationToken string
	RefreshToken      string
}
