package jwt

import "github.com/golang-jwt/jwt/v5"

type claim struct {
	AuthID          int64
	SchoolID        int64
	Role            string
	IsEmailVerified bool
	jwt.RegisteredClaims
}
