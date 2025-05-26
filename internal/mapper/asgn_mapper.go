package mapper

import (
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/models"
)

func UserAndReqAsgnToAsgnResp(user models.Users, req *asgn.AssignmentRequest) asgn.AssignmentResponse {
	return asgn.AssignmentResponse{
		Lecturer: auth.UserSignUpResponse{
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
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
	}
}
