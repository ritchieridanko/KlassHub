package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/clients"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/metadata"
)

type UserHandler struct {
	uc clients.UserClient
}

func NewUserHandler(uc clients.UserClient) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	u, err := h.uc.GetMe(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, true)...,
		),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"User retrieved successfully",
		dtos.UserGetMeResponse{
			User: h.toUser(u),
		},
	)
}

func (h *UserHandler) toUser(u *models.User) *dtos.User {
	if u == nil {
		return nil
	}
	return &dtos.User{
		ID:             u.ID.String(),
		SchoolUserID:   u.SchoolUserID,
		Role:           u.Role,
		Name:           u.Name,
		Nickname:       u.Nickname,
		Birthplace:     u.Birthplace,
		Birthdate:      u.Birthdate,
		Sex:            u.Sex,
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}
