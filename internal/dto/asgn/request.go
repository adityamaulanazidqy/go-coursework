package asgn

import (
	"go-coursework/internal/models"
	"mime/multipart"
	"time"
)

type AssignmentRequest struct {
	ID             int                   `json:"id" gorm:"primary_key;"`
	LecturerID     int                   `json:"lecturer_id"`
	SemesterID     int                   `json:"semester_id"`
	StudyProgramID int                   `json:"study_program_id"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	FileHeader     *multipart.FileHeader `json:"-"`
	MultipartFile  multipart.File        `json:"-"`
	Deadline       time.Time             `json:"deadline"`
	IsActive       bool                  `json:"-"`
}

type AssignmentUpdateRequest struct {
	ID            int                   `json:"id" gorm:"primary_key;"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	FileHeader    *multipart.FileHeader `json:"-"`
	MultipartFile multipart.File        `json:"-"`
	Deadline      time.Time             `json:"deadline"`
	IsActive      bool                  `json:"-"`
	OriginalData  models.Assignment     `json:"-"`
}

type AssignmentCommentRequest struct {
	AssignmentID int    `json:"-"`
	UserID       int    `json:"-"`
	Content      string `json:"content"`
}

type DeleteComment struct {
	UserID       int    `json:"-"`
	AssignmentID int    `json:"-"`
	CommentID    string `json:"-"`
}
