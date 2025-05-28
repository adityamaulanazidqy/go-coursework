package asgn

import (
	"go-coursework/internal/dto/auth"
	"time"
)

type AssignmentResponse struct {
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
