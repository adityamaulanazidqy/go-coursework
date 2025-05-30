package mapper

import (
	"go-coursework/internal/dto/settings"
	"go-coursework/internal/models"
)

func UsersToSetResponse(user *models.Users) settings.SetResponse {
	return settings.SetResponse{
		Username:     user.Username,
		Email:        user.Email,
		Telephone:    user.Telephone,
		StudyProgram: user.StudyProgram.Name,
		Semester:     user.Semester.Name,
		Role:         user.Role.Name,
		Batch:        user.Batch,
		Profile:      user.Profile,
	}
}

func ExistingToUsers(exitingUser *models.Users, req *settings.UpdateUserInfo, filename string) models.Users {
	return models.Users{
		Username:       req.Username,
		Email:          req.Email,
		Telephone:      exitingUser.Telephone,
		StudyProgramID: exitingUser.StudyProgramID,
		Password:       exitingUser.Password,
		RoleID:         exitingUser.RoleID,
		Profile:        filename,
		Batch:          exitingUser.Batch,
	}
}
