package models

import "time"

type Submission struct {
	ID                  int       `json:"id" gorm:"primary_key;"`
	AssignmentID        int       `json:"assignment_id"`
	StudentID           int       `json:"student_id"`
	FileURL             string    `json:"file_url"`
	SubmittedAt         time.Time `json:"submitted_at"`
	StatusSubmissionsID int       `json:"status_submissions_id"`
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
