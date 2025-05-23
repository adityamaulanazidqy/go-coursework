package models

import "go-coursework/constants"

type Users struct {
	ID             int                     `json:"id" gorm:"primary_key;"`
	Username       string                  `json:"username"`
	Email          string                  `json:"email"`
	Telephone      *string                 `json:"telephone"`
	StudyProgramID int                     `json:"study_program_id"`
	StudyProgram   constants.StudyPrograms `json:"study_program" gorm:"foreignkey:StudyProgramID"`
	Password       string                  `json:"password"`
	RoleID         int                     `json:"role_id"`
	Role           constants.Roles         `json:"role" gorm:"foreignkey:RoleID"`
	Profile        *string                 `json:"profile"`
	Batch          int                     `json:"batch"`
}

func (Users) TableName() string {
	return "users"
}
