package dtos

type CreateSchoolProfileResponse struct {
	School      *School      `json:"school,omitempty"`
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type CreateUserAccountResponse struct {
	Auth *AuthAdmin `json:"auth,omitempty"`
	User *UserAdmin `json:"user,omitempty"`
}
