package notification

type SaveFCMTokenReq struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

type WaVerificationReq struct {
	Otp int `json:"otp"`
}
