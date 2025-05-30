package mapper

import (
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/models"
)

func UserAndReqAsgnToAsgnResp(user models.Users, req *asgn.AssignmentRequest, filename string) asgn.AssignmentResponse {
	return asgn.AssignmentResponse{
		Lecturer: auth.UserSignUpResponse{
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
		Lecturer: auth.UserSignUpResponse{
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
