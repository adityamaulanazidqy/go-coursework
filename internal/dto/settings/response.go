package settings

type SetResponse struct {
	ID           int     `json:"-"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Telephone    *string `json:"telephone"`
	StudyProgram string  `json:"study_program"`
	Role         string  `json:"role"`
	Batch        int     `json:"batch"`
	Profile      string  `json:"profile"`
}

type UpdateResponse struct {
	ID           int     `json:"-"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Telephone    *string `json:"telephone"`
	StudyProgram string  `json:"study_program"`
	Role         string  `json:"role"`
	Batch        int     `json:"batch"`
	Profile      string  `json:"profile"`
}
