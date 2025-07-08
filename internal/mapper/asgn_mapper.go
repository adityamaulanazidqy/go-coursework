package mapper

import (
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/models"
	"time"
)

func UserAndReqAsgnToAsgnResp(user models.Users, req *asgn.AssignmentRequest, filename string, id int) asgn.AssignmentResponse {
	return asgn.AssignmentResponse{
		ID: id,
		Lecturer: auth.UserSignUpResponse{
			ID:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			EmailVerified:     user.ContactVerification.EmailVerified,
			Telephone:         user.Telephone,
			TelephoneVerified: user.ContactVerification.TelephoneVerified,
			StudyProgram:      user.StudyProgram.Name,
			Semester:          user.Semester.Name,
			Role:              user.Role.Name,
			Batch:             user.Batch,
			Profile:           user.Profile,
		},
		Title:       req.Title,
		Filename:    filename,
		Description: req.Description,
		Deadline:    req.Deadline,
	}
}

func MapAssignmentToResponse(user *models.Users, assignment *models.Assignment) asgn.AssignmentResponse {
	return asgn.AssignmentResponse{
		ID: assignment.ID,
		Lecturer: auth.UserSignUpResponse{
			ID:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			EmailVerified:     user.ContactVerification.EmailVerified,
			Telephone:         user.Telephone,
			TelephoneVerified: user.ContactVerification.TelephoneVerified,
			StudyProgram:      user.StudyProgram.Name,
			Semester:          user.Semester.Name,
			Role:              user.Role.Name,
			Batch:             user.Batch,
			Profile:           user.Profile,
		},
		Title:       assignment.Title,
		Description: assignment.Description,
		Deadline:    assignment.Deadline,
		Filename:    assignment.Filename,
	}
}

func MapCommentToResponse(comment *models.AssignmentComment, user *models.Users) asgn.CommentResponse {
	return asgn.CommentResponse{
		User: auth.UserSignUpResponse{
			ID:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			EmailVerified:     user.ContactVerification.EmailVerified,
			Telephone:         user.Telephone,
			TelephoneVerified: user.ContactVerification.TelephoneVerified,
			StudyProgram:      user.StudyProgram.Name,
			Role:              user.Role.Name,
			Batch:             user.Batch,
			Profile:           user.Profile,
		},
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func MapSubmissionResponse(user models.Users, lecturer models.Users, assignment models.Assignment, status string, submitted time.Time) asgn.SubmissionResponse {
	return asgn.SubmissionResponse{
		User: auth.UserSignUpResponse{
			ID:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			EmailVerified:     user.ContactVerification.EmailVerified,
			Telephone:         user.Telephone,
			TelephoneVerified: user.ContactVerification.TelephoneVerified,
			StudyProgram:      user.StudyProgram.Name,
			Semester:          user.Semester.Name,
			Role:              user.Role.Name,
			Batch:             user.Batch,
			Profile:           user.Profile,
		},
		Assignment: asgn.AssignmentResponse{
			ID: assignment.ID,
			Lecturer: auth.UserSignUpResponse{
				ID:                lecturer.ID,
				Username:          lecturer.Username,
				Email:             lecturer.Email,
				EmailVerified:     lecturer.ContactVerification.EmailVerified,
				Telephone:         lecturer.Telephone,
				TelephoneVerified: lecturer.ContactVerification.TelephoneVerified,
				StudyProgram:      lecturer.StudyProgram.Name,
				Semester:          lecturer.Semester.Name,
				Role:              lecturer.Role.Name,
				Batch:             lecturer.Batch,
				Profile:           lecturer.Profile,
			},
			Title:       assignment.Title,
			Filename:    assignment.Filename,
			Description: assignment.Description,
			Deadline:    assignment.Deadline,
		},
		Status:      status,
		SubmittedAt: submitted,
	}

}

func MapUserToUserSignUpResponse(user models.Users) *auth.UserSignUpResponse {
	return &auth.UserSignUpResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		EmailVerified:     user.ContactVerification.EmailVerified,
		Telephone:         user.Telephone,
		TelephoneVerified: user.ContactVerification.TelephoneVerified,
		StudyProgram:      user.StudyProgram.Name,
		Semester:          user.Semester.Name,
		Role:              user.Role.Name,
		Batch:             user.Batch,
		Profile:           user.Profile,
	}
}
