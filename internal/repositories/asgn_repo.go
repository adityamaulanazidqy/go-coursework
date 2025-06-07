package repositories

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/constants"
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/helpers"
	"go-coursework/internal/mapper"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
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
		Preload("ContactVerification").
		Preload("Semester").
		Where("id = ?", req.LecturerID).
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
		LecturerID:     req.LecturerID,
		SemesterID:     req.SemesterID,
		StudyProgramID: req.StudyProgramID,
		Title:          req.Title,
		Description:    req.Description,
		Filename:       filename,
		Deadline:       req.Deadline,
		IsActive:       req.IsActive,
	}

	if err := r.rctx.DB.Create(&assignment).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message,
			pkgerr.ErrInternalServer.Details
	}

	resp = mapper.UserAndReqAsgnToAsgnResp(user, req, filename, assignment.ID)

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
		Preload("Semester").
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
		ID:          req.ID,
		Lecturer:    lecturer,
		Title:       req.Title,
		Filename:    req.Filename,
		Description: req.Description,
		Deadline:    req.Deadline,
	}

	return resp, http.StatusOK, op, nil, "Successfully retrieved Assignment", nil
}

func (r *AssignmentRepo) GetAll(userID int) (
	resp []asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {

	const op = "repositories.AssignmentRepo.GetAll"

	var (
		semesterID     int
		studyProgramID int
	)

	var exitingUser models.Users
	if err := r.rctx.DB.Where("id = ?", userID).Find(&exitingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	semesterID = exitingUser.SemesterID
	studyProgramID = exitingUser.StudyProgramID

	var assignments []models.Assignment
	if err := r.rctx.DB.Where("semester_id = ? AND study_program_id = ?", semesterID, studyProgramID).Find(&assignments).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if len(assignments) == 0 {
		err = errors.New("no assignments found")
		return resp, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details
	}

	for _, assignment := range assignments {
		var users []models.Users
		if err := r.rctx.DB.
			Preload("ContactVerification").
			Preload("StudyProgram").
			Preload("Role").
			Preload("Semester").
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
				StudyProgram:      user.StudyProgram.Name,
				Semester:          user.Semester.Name,
				Role:              user.Role.Name,
				Batch:             user.Batch,
				Profile:           user.Profile,
			}

			resp = append(resp, asgn.AssignmentResponse{
				ID:          assignment.ID,
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

func (r *AssignmentRepo) Update(assignment *models.Assignment, req *asgn.AssignmentUpdateRequest) (
	resp asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentRepo.Update"

	hasChanges := false
	updateFields := make(map[string]interface{})

	if req.Title != "" && assignment.Title != req.Title {
		hasChanges = true
		updateFields["title"] = req.Title
	}

	if req.Description != "" && assignment.Description != req.Description {
		hasChanges = true
		updateFields["description"] = req.Description
	}

	if !assignment.Deadline.Equal(req.Deadline) {
		hasChanges = true
		updateFields["deadline"] = req.Deadline
	}

	var oldFilename string
	if req.MultipartFile != nil && req.FileHeader != nil {
		hasChanges = true
		oldFilename = assignment.Filename
	}

	if !hasChanges {
		return resp, http.StatusBadRequest, op, errors.New("no change"),
			pkgerr.ErrAssignmentNotUpdated.Message, pkgerr.ErrAssignmentNotUpdated.Details
	}

	tx := r.rctx.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if req.MultipartFile != nil && req.FileHeader != nil {
		filename, err := helpers.SaveImages().Asgn(req.MultipartFile, req.FileHeader, "_")
		if err != nil {
			return resp, fiber.StatusInternalServerError, op, err,
				pkgerr.ErrFileSave.Message, pkgerr.ErrFileSave.Details
		}
		updateFields["filename"] = filename
		assignment.Filename = filename
	}

	var user models.Users
	if err := tx.
		Preload("StudyProgram").
		Preload("Role").
		Preload("ContactVerification").
		Preload("Semester").
		First(&user, "id = ?", assignment.LecturerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, op, err,
				pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}
		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if err := tx.Model(&models.Assignment{}).
		Where("id = ?", assignment.ID).
		Updates(updateFields).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if oldFilename != "" {
		if err := helpers.DeleteImages().Assignment(oldFilename); err != nil {
			r.rctx.Logger.Logger.Error("failed to delete old file")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
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

func (r *AssignmentRepo) GetAssignmentLecturer(lecturerID int) (
	resp []asgn.AssignmentResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {

	const op = "repositories.AssignmentRepo.GetAssignmentLecturer"

	var assignments []models.Assignment
	if err := r.rctx.DB.Where("lecturer_id = ?", lecturerID).Find(&assignments).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	if len(assignments) == 0 {
		err = errors.New("no assignments found")
		return resp, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details
	}

	for _, assignment := range assignments {
		var users []models.Users
		if err := r.rctx.DB.
			Preload("ContactVerification").
			Preload("StudyProgram").
			Preload("Role").
			Preload("Semester").
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
				StudyProgram:      user.StudyProgram.Name,
				Semester:          user.Semester.Name,
				Role:              user.Role.Name,
				Batch:             user.Batch,
				Profile:           user.Profile,
			}

			resp = append(resp, asgn.AssignmentResponse{
				ID:          assignment.ID,
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

func (r *AssignmentRepo) Submissions(submission *models.Submission, assignment models.Assignment) (
	resp asgn.SubmissionResponse,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentsRepo.Submissions"

	var (
		wg             sync.WaitGroup
		userResult     models.Users
		lecturerResult models.Users
		errOnce        sync.Once
		firstErr       error
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		user, c, o, e, m, d := r.searchUserSignUpResponse(submission.StudentID)
		if e != nil {
			errOnce.Do(func() {
				firstErr = e
				code = c
				opRepo = o
				msg = m
				details = d
			})
			return
		}
		userResult = user
	}()

	go func() {
		defer wg.Done()
		lecturer, c, o, e, m, d := r.searchUserSignUpResponse(assignment.LecturerID)
		if e != nil {
			errOnce.Do(func() {
				firstErr = e
				code = c
				opRepo = o
				msg = m
				details = d
			})
			return
		}
		lecturerResult = lecturer
	}()

	wg.Wait()

	if firstErr != nil {
		return resp, code, opRepo, firstErr, msg, details
	}

	var existingSubmissions models.Submission
	if err := r.rctx.DB.Where("assignment_id = ? AND student_id = ?", submission.AssignmentID, submission.StudentID).First(&existingSubmissions).Error; err == nil {
		err = errors.New("submission already exists")
		return resp, fiber.StatusConflict, op, err, pkgerr.ErrSubmissionAlreadyExists.Message, pkgerr.ErrSubmissionAlreadyExists.Details
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.rctx.DB.WithContext(ctx).Create(&submission).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, pkgerr.ErrSubmissionFailed.Message, pkgerr.ErrSubmissionFailed.Details
	}

	var status constants.StatusSubmissions
	if err := r.rctx.DB.Where("id = ?", submission.StatusSubmissionsID).First(&status).Error; err != nil {
		return asgn.SubmissionResponse{}, fiber.StatusInternalServerError, op, err, pkgerr.ErrSubmissionStatusNotFound.Message, pkgerr.ErrSubmissionStatusNotFound.Details
	}

	resp = mapper.MapSubmissionResponse(userResult, lecturerResult, assignment, status.Name, submission.SubmittedAt)

	return resp, fiber.StatusCreated, op, nil, "Successfully Submission", nil
}

func (r *AssignmentRepo) searchUserSignUpResponse(id int) (
	user models.Users,
	code int,
	opRepo string,
	err error,
	msg string,
	details []string,
) {
	const op = "repositories.AssignmentsRepo.searchUserSignUpResponse"

	if err := r.rctx.DB.
		Preload("StudyProgram").
		Preload("Semester").
		Preload("Role").
		Preload("ContactVerification").
		Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details
		}

		return user, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details
	}

	return user, fiber.StatusOK, op, nil, "Successfully searching data user", nil
}
