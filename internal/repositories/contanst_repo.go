package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/constants"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"gorm.io/gorm"
)

type ConstRepo struct {
	rctx *models.RouterContext
}

func NewConstRepo(rctx *models.RouterContext) *ConstRepo {
	return &ConstRepo{
		rctx: rctx,
	}
}

func (r *ConstRepo) GetSemesters() (resp []constants.Semesters, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.GetSemesters"

	var semesters []constants.Semesters
	if err := r.rctx.DB.Find(&semesters).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return semesters, fiber.StatusOK, op, nil, "Successfully retrieved the semesters", nil
}

func (r *ConstRepo) PostSemesters(req *constants.Semesters) (resp *constants.Semesters, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.PostSemesters"

	if err := r.rctx.DB.Create(&req).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return req, fiber.StatusOK, op, nil, "Successfully created the semesters", nil
}

func (r *ConstRepo) DeleteSemester(semesterID int) (resp constants.Semesters, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.DeleteSemester"

	if err := r.rctx.DB.Find(&resp, "id = ?", semesterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details
		}

		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var semester constants.Semesters
	if err := r.rctx.DB.Delete(&semester, "id = ?", semesterID).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return resp, fiber.StatusOK, op, nil, "Successfully deleted the semester", nil
}

func (r *ConstRepo) GetStudyPrograms() (resp []constants.StudyPrograms, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.GetStudyPrograms"

	var studyPrograms []constants.StudyPrograms
	if err := r.rctx.DB.Find(&studyPrograms).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return studyPrograms, fiber.StatusOK, op, nil, "Successfully retrieved the study programs", nil
}

func (r *ConstRepo) PostStudyPrograms(req *constants.StudyPrograms) (resp *constants.StudyPrograms, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.PostStudyPrograms"

	if err := r.rctx.DB.Create(&req).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return req, fiber.StatusOK, op, nil, "Successfully created the study programs", nil
}

func (r *ConstRepo) DeleteProgram(req int) (resp constants.StudyPrograms, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.constRepo.DeleteProgram"

	if err := r.rctx.DB.Find(&resp, "id = ?", req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details
		}

		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	var studyProgram constants.StudyPrograms
	if err := r.rctx.DB.Delete(&studyProgram, "id = ?", req).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return resp, fiber.StatusOK, op, nil, "Successfully deleted the study programs", nil
}
