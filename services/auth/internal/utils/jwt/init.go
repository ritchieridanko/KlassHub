package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
}

func Init(issuer, secret string, dn time.Duration) *JWT {
	return &JWT{issuer: issuer, secret: secret, duration: dn}
}

func (j *JWT) Generate(authID, schoolID int64, role string, isEmailVerified bool, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now().UTC()
		now = &t
	}
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim{
			AuthID:          authID,
			SchoolID:        schoolID,
			Role:            role,
			IsEmailVerified: isEmailVerified,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   fmt.Sprintf("auth_%d", authID),
				IssuedAt:  &jwt.NumericDate{Time: *now},
				ExpiresAt: &jwt.NumericDate{Time: now.Add(j.duration)},
			},
		},
	).SignedString([]byte(j.secret))
}
