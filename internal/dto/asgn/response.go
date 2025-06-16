package asgn

import (
	"go-coursework/internal/dto/auth"
	"time"
)

type AssignmentResponse struct {
	ID          int                     `json:"id"`
	Lecturer    auth.UserSignUpResponse `json:"lecturer"`
	Title       string                  `json:"title"`
	Filename    string                  `json:"filename"`
	Description string                  `json:"description"`
	Deadline    time.Time               `json:"deadline"`
}

type CommentResponse struct {
	User      auth.UserSignUpResponse `json:"user"`
	Content   string                  `json:"content"`
	CreatedAt time.Time               `json:"created_at"`
}

type SubmissionResponse struct {
	User        auth.UserSignUpResponse `json:"user"`
	Assignment  AssignmentResponse      `json:"assignment"`
	Status      string                  `json:"status"`
	SubmittedAt time.Time               `json:"submitted_at"`
}

type GetSubmissionsResponse struct {
	ID          int                     `json:"id"`
	User        auth.UserSignUpResponse `json:"user"`
	Status      string                  `json:"status"`
	SubmittedAt time.Time               `json:"submitted_at"`
}

type SubmissionGradeResponse struct {
	Submission SubmissionResponse      `json:"submission"`
	Lecturer   auth.UserSignUpResponse `json:"lecturer"`
}

type SubmissionUpdateResponse struct {
	User       auth.UserSignUpResponse `json:"user"`
	Assignment AssignmentResponse      `json:"assignment"`
	Status     string                  `json:"status"`
	UpdatedAt  time.Time               `json:"updated_at"`
}
