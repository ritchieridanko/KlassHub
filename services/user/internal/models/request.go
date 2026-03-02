package models

type GetUserRequest struct {
	AuthID   int64
	SchoolID int64
}

type GetUserAuthInfoRequest struct {
	AuthID int64
}
