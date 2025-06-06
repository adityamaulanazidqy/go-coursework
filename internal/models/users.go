package models

import "go-coursework/constants"

type Users struct {
	ID                  int                     `json:"id" gorm:"primary_key;"`
	Username            string                  `json:"username"`
	Email               string                  `json:"email"`
	Telephone           *string                 `json:"telephone"`
	StudyProgramID      int                     `json:"study_program_id"`
	StudyProgram        constants.StudyPrograms `json:"study_program" gorm:"foreignkey:StudyProgramID"`
	SemesterID          int                     `json:"semester_id"`
	Semester            constants.Semesters     `json:"semester" gorm:"foreignkey:SemesterID"`
	Password            string                  `json:"password"`
	RoleID              int                     `json:"role_id"`
	Role                constants.Roles         `json:"role" gorm:"foreignkey:RoleID"`
	Profile             string                  `json:"profile" gorm:"default:icon_default.jpg"`
	Batch               int                     `json:"batch"`
	ContactVerification UserContactVerification `json:"contact_verification" gorm:"foreignkey:ID"`
}

func (Users) TableName() string {
	return "users"
}
