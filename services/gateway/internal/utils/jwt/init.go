package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
)

type JWT struct {
	secret string
}

func Init(secret string) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) Parse(token string) (*claim, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&claim{},
		func(t *jwt.Token) (any, error) {
			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := t.Claims.(*claim)
	if !ok {
		return nil, ce.ErrInvalidJWTClaim
	}
	return claim, nil
}
