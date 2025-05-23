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
		Role:         user.Role.Name,
		Batch:        user.Batch,
		Profile:      user.Profile,
	}
}
