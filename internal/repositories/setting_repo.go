package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/dto/settings"
	"go-coursework/internal/helpers"
	"go-coursework/internal/logger"
	"go-coursework/internal/mapper"
	"go-coursework/internal/models"
	"gorm.io/gorm"
)

type SettingRepo struct {
	db        *gorm.DB
	logLogrus *logger.ErrorLogger
}

func NewSettingRepo(db *gorm.DB) *SettingRepo {
	return &SettingRepo{db: db}
}

func (r *SettingRepo) SetProfile(req *settings.SetProfile, userID int) (resp settings.SetResponse, code int, opRepo string, err error, msg string, details []string) {
	opRepo = "repositories.Setting.SetProfile"

	var (
		msgNotFoundUser               = "User Not Found"
		msgNotFoundUserDetail         = []string{"Please enter the user ID correctly"}
		msgFailedSetProfile           = "Failed to set profile"
		msgFailedSetProfileDetail     = []string{"Make sure your photo format meets the requirements", "Make sure your internet is on"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	filename, err := helpers.SaveImages().Profile(req.MultipartFile, req.FileHeader, "_")
	if err != nil {
		return resp, fiber.StatusInternalServerError, opRepo, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	var user models.Users
	if err := r.db.Model(&user).Where("id = ?", userID).Update("profile", filename).Error; err != nil {
		return resp, fiber.StatusInternalServerError, opRepo, err, msgFailedSetProfile, msgFailedSetProfileDetail
	}

	if err := r.db.
		Preload("StudyProgram").
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, opRepo, err, msgNotFoundUser, msgNotFoundUserDetail
		}

		return resp, fiber.StatusInternalServerError, opRepo, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	resp = mapper.UsersToSetResponse(&user)

	return resp, fiber.StatusCreated, opRepo, nil, "Successfully Set Profile", nil
}

func (r *SettingRepo) SetTelephone(req *settings.SetTelephone, userID int) (resp settings.SetResponse, code int, opRepo string, err error, msg string, details []string) {
	opRepo = "repositories.Setting.SetTelephone"

	var (
		msgNotFoundUser               = "User Not Found"
		msgNotFoundUserDetail         = []string{"Please enter the user ID correctly"}
		msgFailedSetTelephone         = "Failed to set telephone"
		msgFailedSetTelephoneDetail   = []string{"Make sure your telephone format meets the requirements", "Make sure your internet is on"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	var user models.Users
	if err := r.db.Model(&user).Where("id = ?", userID).Update("telephone", req.Telephone).Error; err != nil {
		return resp, fiber.StatusInternalServerError, opRepo, err, msgFailedSetTelephone, msgFailedSetTelephoneDetail
	}

	if err := r.db.
		Preload("StudyProgram").
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, opRepo, err, msgNotFoundUser, msgNotFoundUserDetail
		}

		return resp, fiber.StatusInternalServerError, opRepo, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	resp = mapper.UsersToSetResponse(&user)

	return resp, fiber.StatusCreated, opRepo, nil, "Successfully Set Telephone", nil
}

func (r *SettingRepo) UpdateUserInfo(req *settings.UpdateUserInfo, userID int) (resp settings.UpdateResponse, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.Setting.UpdateUserInfo"

	var (
		msgNotFoundUser               = "User Not Found"
		msgNotFoundUserDetail         = []string{"Please enter the user ID correctly"}
		msgFailedDeleteProfile        = "Failed to delete profile"
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	var exitingUser models.Users
	if err := r.db.Where("id = ?", userID).First(&exitingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, msgNotFoundUser, msgNotFoundUserDetail
		}

		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	if req.Username == "" {
		req.Username = exitingUser.Username
	}

	if req.Email == "" {
		req.Email = exitingUser.Email
	}

	if req.Telephone == "" && exitingUser.Telephone != nil {
		req.Telephone = *exitingUser.Telephone
	}

	var filename string
	if req.Profile.MultipartFile != nil && req.Profile.FileHeader != nil {
		filename, err = helpers.SaveImages().Profile(req.Profile.MultipartFile, req.Profile.FileHeader, "_")
		if err != nil {
			return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
		}

		if exitingUser.Profile != "" {
			err = helpers.DeleteImages().Profile(exitingUser.Profile)
			if err != nil {
				r.logLogrus.LogUserError(exitingUser.Email, err, msgFailedDeleteProfile)
			}
		}
	} else {
		if exitingUser.Profile != "" {
			filename = exitingUser.Profile
		}
	}

	var user = mapper.ExistingToUsers(&exitingUser, req, filename)

	if err := r.db.Model(&user).Where("id = ?", userID).Updates(user).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	if err := r.db.
		Preload("StudyProgram").
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, msgNotFoundUser, msgNotFoundUserDetail
		}

		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	return resp, fiber.StatusCreated, opRepo, nil, "Successfully Update User", nil
}
