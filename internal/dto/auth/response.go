package auth

type UserSignUpResponse struct {
	ID                int     `json:"-"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	EmailVerified     bool    `json:"email_verified"`
	Telephone         *string `json:"telephone"`
	TelephoneVerified bool    `json:"telephone_verified"`
	StudyProgram      string  `json:"study_program"`
	Password          string  `json:"-"`
	Role              string  `json:"role"`
	Batch             int     `json:"batch"`
	Profile           *string `json:"profile"`
}

type UserSignInResponse struct {
	Email             string  `json:"email"`
	EmailVerified     bool    `json:"email_verified"`
	Telephone         *string `json:"telephone"`
	TelephoneVerified bool    `json:"telephone_verified"`
	Role              string  `json:"role"`
	Token             string  `json:"token"`
}
