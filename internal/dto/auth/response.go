package auth

type UserSignUpResponse struct {
	ID           int     `json:"-"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Telephone    *string `json:"telephone"`
	StudyProgram string  `json:"study_program"`
	Password     string  `json:"-"`
	Role         string  `json:"role"`
	Batch        int     `json:"batch"`
	Profile      *string `json:"profile"`
}

type UserSignInResponse struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}
