package notification

import "go-coursework/internal/dto/auth"

type SaveFCMTokenResp struct {
	User  auth.UserSignUpResponse `json:"user"`
	Token string                  `json:"token"`
}
