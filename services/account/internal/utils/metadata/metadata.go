package metadata

import (
	"strconv"

	"github.com/ritchieridanko/klasshub/services/account/internal/constants"
	"github.com/ritchieridanko/klasshub/services/account/internal/models"
)

func Auth(ac *models.AuthContext, auth, school, role, verification bool) []Pair {
	pairs := make([]Pair, 0, 4)
	if auth {
		pairs = append(pairs, NewPair(constants.MDKeyAuthID, strconv.FormatInt(ac.AuthID, 10)))
	}
	if school {
		pairs = append(pairs, NewPair(constants.MDKeySchoolID, strconv.FormatInt(ac.SchoolID, 10)))
	}
	if role {
		pairs = append(pairs, NewPair(constants.MDKeyRole, ac.Role))
	}
	if verification {
		pairs = append(pairs, NewPair(constants.MDKeyIsVerified, strconv.FormatBool(ac.IsVerified)))
	}
	return pairs
}
