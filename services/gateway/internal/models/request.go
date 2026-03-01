package models

type LoginRequest struct {
	Identifier string
	Password   string
	Subdomain  string
}
