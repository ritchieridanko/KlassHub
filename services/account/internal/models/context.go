package models

type AuthContext struct {
	AuthID     int64
	SchoolID   int64
	Role       string
	IsVerified bool
}
