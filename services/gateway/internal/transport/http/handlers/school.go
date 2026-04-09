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

type SchoolHandler struct {
	sc clients.SchoolClient
}

func NewSchoolHandler(sc clients.SchoolClient) *SchoolHandler {
	return &SchoolHandler{sc: sc}
}

func (h *SchoolHandler) GetMe(ctx *gin.Context) {
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

	s, err := h.sc.GetMe(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, false)...,
		),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"School retrieved successfully",
		dtos.SchoolGetMeResponse{
			School: h.toSchool(s),
		},
	)
}

func (h *SchoolHandler) toSchool(s *models.School) *dtos.School {
	if s == nil {
		return nil
	}
	return &dtos.School{
		NPSN:           s.NPSN,
		NPSNVerifiedAt: s.NPSNVerifiedAt,
		Name:           s.Name,
		Level:          s.Level,
		Ownership:      s.Ownership,
		ProfilePicture: s.ProfilePicture,
		ProfileBanner:  s.ProfileBanner,
		Accreditation:  s.Accreditation,
		EstablishedAt:  s.EstablishedAt,
		Province:       s.Province,
		CityRegency:    s.CityRegency,
		District:       s.District,
		Subdistrict:    s.Subdistrict,
		Street:         s.Street,
		Postcode:       s.Postcode,
		Phone:          s.Phone,
		Email:          s.Email,
		Website:        s.Website,
		Timezone:       s.Timezone,
		CreatedAt:      s.CreatedAt,
	}
}
