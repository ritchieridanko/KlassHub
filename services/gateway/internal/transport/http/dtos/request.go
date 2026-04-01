package dtos

type CreateSchoolAuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type VerifyEmailRequest struct {
	VerificationToken string `form:"token" binding:"required"`
}
