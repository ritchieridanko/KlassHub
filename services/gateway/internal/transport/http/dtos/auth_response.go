package dtos

type ChangePasswordResponse struct {
	Auth *Auth `json:"auth,omitempty"`
}

type CreateSchoolAuthResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type EmailAvailabilityCheckResponse struct {
	IsAvailable bool `json:"is_available"`
}

type LoginResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type ResendVerificationResponse struct {
	Email string `json:"email"`
}

type RotateAuthTokenResponse struct {
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type UsernameAvailabilityCheckResponse struct {
	IsAvailable bool `json:"is_available"`
}

type VerifyEmailResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}
