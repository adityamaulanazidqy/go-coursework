package models

import "time"

type Assignment struct {
	ID             int       `json:"id" gorm:"primary_key;"`
	LecturerID     int       `json:"lecturer_id"`
	SemesterID     int       `json:"semester_id"`
	StudyProgramID int       `json:"study_program_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Filename       string    `json:"filename"`
	Deadline       time.Time `json:"deadline"`
	IsActive       bool      `json:"is_active"`
}

func (Assignment) TableName() string {
	return "assignments"
}

type AssignmentComment struct {
	ID           int       `json:"id"`
	AssignmentID int       `json:"assignment_id"`
	UserID       int       `json:"user_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (AssignmentComment) TableName() string {
	return "assignment_comments"
}

type AssignmentFile struct {
	ID           int       `json:"id"`
	AssignmentID int       `json:"assignment_id"`
	FileURL      string    `json:"file_url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

func (AssignmentFile) TableName() string {
	return "assignment_files"
}
