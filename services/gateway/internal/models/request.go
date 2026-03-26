package models

type CreateSchoolAuthReq struct {
	Email    string
	Password string
}

type LoginReq struct {
	Identifier string
	Password   string
}
