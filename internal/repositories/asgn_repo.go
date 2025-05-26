package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/helpers"
	"go-coursework/internal/mapper"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type AssignmentRepo struct {
	rctx *models.RouterContext
}

func NewAssignmentRepo(rctx *models.RouterContext) *AssignmentRepo {
	return &AssignmentRepo{rctx: rctx}
}

func (r *AssignmentRepo) Post(req *asgn.AssignmentRequest) (
	resp asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.Post"

	var user models.Users

	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").Where("id = ?", req.LecturerID).
		First(&user, req.LecturerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err,
				pkgerr.ErrUserNotFound.Message,
				pkgerr.ErrUserNotFound.Details
		}

		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message,
			pkgerr.ErrInternalServer.Details
	}

	filename, err := helpers.SaveImages().Asgn(req.MultipartFile, req.FileHeader, "_")
	if err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var assignment = models.Assignment{
		LecturerID:  req.LecturerID,
		Title:       req.Title,
		Description: req.Description,
		Filename:    filename,
		Deadline:    req.Deadline,
		IsActive:    req.IsActive,
	}

	if err := r.rctx.DB.Create(&assignment).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message,
			pkgerr.ErrInternalServer.Details
	}

	//resp = asgn.AssignmentResponse{
	//	Lecturer: auth.UserSignUpResponse{
	//		Username:          user.Username,
	//		Email:             user.Email,
	//		EmailVerified:     user.ContactVerification.EmailVerified,
	//		Telephone:         user.Telephone,
	//		TelephoneVerified: user.ContactVerification.TelephoneVerified,
	//		StudyProgram:      user.StudyProgram.Name,
	//		Role:              user.Role.Name,
	//		Batch:             user.Batch,
	//		Profile:           user.Profile,
	//	},
	//	Title:       req.Title,
	//	Description: req.Description,
	//	Deadline:    req.Deadline,
	//}

	resp = mapper.UserAndReqAsgnToAsgnResp(user, req)

	return resp, http.StatusOK, op, nil, "Successfully created Assignment", nil
}
