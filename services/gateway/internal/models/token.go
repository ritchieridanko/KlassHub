package models

type AccessToken struct {
	Token     string
	ExpiresIn int64 // seconds
}

type RefreshToken struct {
	Token     string
	ExpiresIn int64 // seconds
}

type AuthToken struct {
	AccessToken  *AccessToken
	RefreshToken *RefreshToken
}
