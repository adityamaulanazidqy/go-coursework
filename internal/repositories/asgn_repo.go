package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/dto/auth"
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

	resp = mapper.UserAndReqAsgnToAsgnResp(user, req, filename)

	return resp, http.StatusOK, op, nil, "Successfully created Assignment", nil
}

func (r *AssignmentRepo) Get(req *models.Assignment) (
	resp asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {

	const op = "repositories.AssignmentRepo.Get"

	var user models.Users
	if err := r.rctx.DB.
		Preload("ContactVerification").
		Preload("StudyProgram").
		Preload("Role").
		First(&user, "id = ?", req.LecturerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}

		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var lecturer = auth.UserSignUpResponse{
		Username:          user.Username,
		Email:             user.Email,
		EmailVerified:     user.ContactVerification.EmailVerified,
		Telephone:         user.Telephone,
		TelephoneVerified: user.ContactVerification.TelephoneVerified,
		StudyProgram:      user.StudyProgram.Name,
		Role:              user.Role.Name,
		Batch:             user.Batch,
		Profile:           user.Profile,
	}

	resp = asgn.AssignmentResponse{
		Lecturer:    lecturer,
		Title:       req.Title,
		Filename:    req.Filename,
		Description: req.Description,
		Deadline:    req.Deadline,
	}

	return resp, http.StatusOK, op, nil, "Successfully retrieved Assignment", nil
}

func (r *AssignmentRepo) GetAll() (
	resp []asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {

	const op = "repositories.AssignmentRepo.GetAll"

	var assignments []models.Assignment
	if err := r.rctx.DB.Find(&assignments).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	for _, assignment := range assignments {
		var users []models.Users
		if err := r.rctx.DB.
			Preload("ContactVerification").
			Preload("StudyProgram").
			Preload("Role").
			Find(&users, "id = ?", assignment.LecturerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
			}

			return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
		}

		for _, user := range users {
			var lecturer = auth.UserSignUpResponse{
				Username:          user.Username,
				Email:             user.Email,
				EmailVerified:     user.ContactVerification.EmailVerified,
				Telephone:         user.Telephone,
				TelephoneVerified: user.ContactVerification.TelephoneVerified,
				Role:              user.Role.Name,
				Batch:             user.Batch,
				Profile:           user.Profile,
			}

			resp = append(resp, asgn.AssignmentResponse{
				Lecturer:    lecturer,
				Title:       assignment.Title,
				Filename:    assignment.Filename,
				Description: assignment.Description,
				Deadline:    assignment.Deadline,
			})
		}
	}

	return resp, http.StatusOK, op, nil, "Successfully retrieved Assignments", nil
}

func (r *AssignmentRepo) Update(assignment *models.Assignment, req *asgn.AssignmentRequest) (
	resp asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.Update"

	var user models.Users
	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").
		First(&user, "id = ?", assignment.LecturerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if req.MultipartFile != nil && req.FileHeader != nil {
		filename, err := helpers.SaveImages().Asgn(req.MultipartFile, req.FileHeader, "_")
		if err != nil {
			return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrFileSave.Message, pkgerr.ErrFileSave.Details
		}
		if assignment.Filename != filename {
			assignment.Filename = filename
		}
	}

	if assignment.Title != req.Title {
		assignment.Title = req.Title
	}
	if assignment.Description != req.Description {
		assignment.Description = req.Description
	}
	if assignment.Deadline != req.Deadline {
		assignment.Deadline = req.Deadline
	}

	noChange := assignment.Title == req.Title &&
		assignment.Description == req.Description &&
		assignment.Deadline.Equal(req.Deadline) &&
		(req.FileHeader == nil || assignment.Filename == assignment.Filename)

	if noChange {
		return resp, http.StatusBadRequest, op, errors.New("no change"),
			pkgerr.ErrAssignmentNotUpdated.Message, pkgerr.ErrAssignmentNotUpdated.Details
	}

	tx := r.rctx.DB.Model(&assignment).Where("id = ?", assignment.ID).Updates(assignment)
	if tx.Error != nil {
		return resp, fiber.StatusInternalServerError, op, tx.Error,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if tx.RowsAffected == 0 {
		return resp, fiber.StatusNotFound, op, nil,
			pkgerr.ErrAssignmentNotUpdated.Message, pkgerr.ErrAssignmentNotUpdated.Details
	}

	resp = mapper.MapAssignmentToResponse(&user, assignment)

	return resp, http.StatusOK, op, nil, "Successfully updated assignment", nil
}

func (r *AssignmentRepo) Delete(req *models.Assignment) (
	resp asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.Delete"

	var user models.Users
	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").
		First(&user, "id = ?", req.LecturerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if err := r.rctx.DB.Model(&req).Where("id = ?", req.ID).Delete(&req).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrDeleteFailed.Message, pkgerr.ErrDeleteFailed.Details
	}

	resp = mapper.MapAssignmentToResponse(&user, req)

	return resp, http.StatusOK, op, nil, "Successfully deleted assignment", nil
}

func (r *AssignmentRepo) Comment(req *asgn.AssignmentCommentRequest) (
	resp asgn.CommentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.Comment"

	var user models.Users
	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").
		First(&user, "id = ?", req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var comment = models.AssignmentComment{
		AssignmentID: req.AssignmentID,
		UserID:       user.ID,
		Content:      req.Content,
	}

	if err := r.rctx.DB.Create(&comment).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrCommentSave.Message, pkgerr.ErrCommentSave.Details
	}

	resp = mapper.MapCommentToResponse(&comment, &user)

	return resp, fiber.StatusOK, op, nil, "Successfully created comment", nil
}

func (r *AssignmentRepo) GetCommentsByAssignmentID(assignmentID int) (
	resp []asgn.CommentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.GetCommentsByAssignmentID"

	var comments []models.AssignmentComment
	if err := r.rctx.DB.
		Where("assignment_id = ?", assignmentID).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {

		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	for _, comment := range comments {
		var user models.Users
		if err := r.rctx.DB.
			Preload("StudyProgram").
			Preload("Role").
			Preload("ContactVerification").
			First(&user, "id = ?", comment.UserID).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, fiber.StatusNotFound, op, err,
					pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
			}

			return resp, fiber.StatusInternalServerError, op, err,
				pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
		}

		resp = append(resp, mapper.MapCommentToResponse(&comment, &user))
	}

	if resp == nil {
		return resp, fiber.StatusNotFound, op, err, pkgerr.ErrNoComments.Message, pkgerr.ErrNoComments.Details
	}

	return resp, fiber.StatusOK, op, nil, "Successfully fetched comments", nil
}

func (r *AssignmentRepo) DeleteComment(req *asgn.DeleteComment) (
	resp asgn.CommentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.DeleteComment"

	var user models.Users
	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").
		First(&user, "id = ?", req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var comment models.AssignmentComment
	if err := r.rctx.DB.First(&comment, "id = ?", req.CommentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrCommentNotFound.Message, pkgerr.ErrCommentNotFound.Details
		}
	}

	if err := r.rctx.DB.Model(&comment).Where("id = ?", req.CommentID).Delete(&comment).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrDeleteComment.Message, pkgerr.ErrDeleteComment.Details
	}

	resp = mapper.MapCommentToResponse(&comment, &user)

	return resp, fiber.StatusOK, op, nil, "Successfully deleted comment", nil
}
