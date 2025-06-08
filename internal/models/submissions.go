package models

import "time"

type Submission struct {
	ID                  int       `json:"id" gorm:"primary_key;"`
	AssignmentID        int       `json:"assignment_id"`
	StudentID           int       `json:"student_id"`
	FileURL             string    `json:"file_url"`
	SubmittedAt         time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	StatusSubmissionsID int       `json:"status_submissions_id" gorm:"default:1"`
}

func (Submission) TableName() string {
	return "submissions"
}

type SubmissionHistories struct {
	ID                  int       `json:"id" gorm:"primary_key;"`
	SubmissionID        int       `json:"submission_id"`
	FileURL             string    `json:"file_url"`
	StatusSubmissionsID int       `json:"status_submissions_id"`
	ChangedAt           time.Time `json:"changed_at"`
	Notes               *string   `json:"notes"`
}

func (SubmissionHistories) TableName() string {
	return "submission_histories"
}

type SubmissionGrades struct {
	ID           int       `json:"id" gorm:"primary_key;"`
	SubmissionID int       `json:"submission_id"`
	LecturerID   int       `json:"lecturer_id"`
	Grade        int       `json:"grade"`
	Notes        *string   `json:"notes"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (SubmissionGrades) TableName() string {
	return "submission_grades"
}
