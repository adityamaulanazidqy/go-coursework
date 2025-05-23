package auth

import "go-coursework/constants"

type SignUpRequest struct {
	ID             int                     `json:"id" gorm:"primary_key;"`
	Username       string                  `json:"username"`
	Email          string                  `json:"email"`
	Telephone      *string                 `json:"telephone"`
	StudyProgramID int                     `json:"study_program_id"`
	StudyProgram   constants.StudyPrograms `json:"-" gorm:"foreignkey:StudyProgramID"`
	Password       string                  `json:"password"`
	RoleID         int                     `json:"role_id"`
	Role           constants.Roles         `json:"-" gorm:"foreignkey:RoleID"`
	Batch          int                     `json:"batch"`
	Profile        *string                 `json:"profile"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
