package dtos

type UserGetMeResponse struct {
	User *User `json:"user,omitempty"`
}
