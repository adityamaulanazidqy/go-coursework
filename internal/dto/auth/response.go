package auth

type UserSignUpResponse struct {
	ID                int     `json:"-"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	EmailVerified     bool    `json:"email_verified"`
	Telephone         *string `json:"telephone"`
	TelephoneVerified bool    `json:"telephone_verified"`
	StudyProgram      string  `json:"study_program"`
	Semester          string  `json:"semester"`
	Password          string  `json:"-"`
	Role              string  `json:"role"`
	Batch             int     `json:"batch"`
	Profile           string  `json:"profile"`
}

type UserSignInResponse struct {
	Email             string  `json:"email"`
	EmailVerified     bool    `json:"email_verified"`
	Telephone         *string `json:"telephone"`
	TelephoneVerified bool    `json:"telephone_verified"`
	Role              string  `json:"role"`
	Semester          string  `json:"semester"`
	Token             string  `json:"token"`
}
