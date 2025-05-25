package otp

type SendOtpEmail struct {
	Email string `json:"email" binding:"required"`
}

type VerifyOtpEmail struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}
