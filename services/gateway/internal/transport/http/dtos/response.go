package dtos

import "time"

type Response[T any] struct {
	Status   int               `json:"status"`
	Message  string            `json:"message"`
	Data     T                 `json:"data,omitempty"`
	Metadata *ResponseMetadata `json:"metadata,omitempty"`
}

type ResponseMetadata struct {
	RequestID string    `json:"request_id"`
	Page      *int      `json:"page,omitempty"`
	PageSize  *int      `json:"page_size,omitempty"`
	Total     *int      `json:"total,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type CreateSchoolAuthResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type LoginResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}

type VerifyEmailResponse struct {
	Auth        *Auth        `json:"auth,omitempty"`
	AccessToken *AccessToken `json:"access_token,omitempty"`
}
