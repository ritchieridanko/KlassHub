package models

type CreateSchoolAuthReq struct {
	SchoolID int64
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
