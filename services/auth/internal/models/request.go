package models

type RequestMeta struct {
	UserAgent string
	IPAddress string
}

type CreateSessionRequest struct {
	AuthID          int64
	SchoolID        int64
	Role            string
	IsEmailVerified bool
	RequestMeta
}

type LoginRequest struct {
	Identifier string
	Password   string
	Subdomain  string
	RequestMeta
}
