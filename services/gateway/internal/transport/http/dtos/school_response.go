package dtos

type CreateSchoolProfileResponse struct {
	School      *School      `json:"school,omitempty"`
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type SchoolGetMeResponse struct {
	School *School `json:"school,omitempty"`
}
