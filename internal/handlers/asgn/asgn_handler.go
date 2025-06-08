package asgn

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/models"
	"go-coursework/internal/repositories"
	pkgerr "go-coursework/pkg/errors"
	"go-coursework/pkg/jwt"
	"mime/multipart"
	"strconv"
	"time"
)

type AssignmentsHandler struct {
	rctx     *models.RouterContext
	asgnRepo *repositories.AssignmentRepo
}

func NewAssignmentsHandler(rctx *models.RouterContext) *AssignmentsHandler {
	return &AssignmentsHandler{rctx: rctx, asgnRepo: repositories.NewAssignmentRepo(rctx)}
}

func (h *AssignmentsHandler) Post(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Post"

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err,
			pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	multipartFile, err := fileHeader.Open()
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err,
			pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	deadlineStr := ctx.FormValue("deadline")
	deadline, err := time.Parse("2006-01-02 15:04:05", deadlineStr)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err,
			pkgerr.ErrInvalidDeadline.Message, pkgerr.ErrInvalidDeadline.Details)
	}

	semesterID, err := strconv.Atoi(ctx.FormValue("semester_id"))
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	studyProgramID, err := strconv.Atoi(ctx.FormValue("study_program_id"))
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	req := asgn.AssignmentRequest{
		LecturerID:     claims.UserID,
		SemesterID:     semesterID,
		StudyProgramID: studyProgramID,
		Title:          ctx.FormValue("title"),
		Description:    ctx.FormValue("description"),
		FileHeader:     fileHeader,
		MultipartFile:  multipartFile,
		Deadline:       deadline,
		IsActive:       true,
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.Post(&req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) Get(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Get"

	claimsUser, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claimsUser == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err,
			pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	claims, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.Get(claims)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) GetAll(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.GetAll"

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err,
			pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetAll(claims.UserID)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) Update(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Update"

	claimsUser, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claimsUser == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(claimsUser.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if assignment.LecturerID != claimsUser.UserID {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, errors.New("unauthorized: lecturer ID mismatch"), pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	fileHeader, err := ctx.FormFile("file")
	var multipartFile multipart.File

	if err == nil && fileHeader != nil {
		multipartFile, err = fileHeader.Open()
		if err != nil {
			return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrFileOpen.Message, pkgerr.ErrFileOpen.Details)
		}
		defer multipartFile.Close()
	}

	deadlineStr := ctx.FormValue("deadline")
	var deadline time.Time

	if deadlineStr != "" {
		deadline, err = time.Parse("2006-01-02 15:04:05", deadlineStr)
		if err != nil {
			return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, "Invalid deadline format. Use YYYY-MM-DD HH:MM:SS", nil)
		}
	} else {
		deadline = assignment.Deadline
	}

	var req = asgn.AssignmentUpdateRequest{
		Title:         ctx.FormValue("title"),
		Description:   ctx.FormValue("description"),
		FileHeader:    fileHeader,
		MultipartFile: multipartFile,
		Deadline:      deadline,
		OriginalData:  *assignment,
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.Update(assignment, &req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) Delete(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Delete"

	claimsUser, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claimsUser == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	claims, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(claimsUser.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if claims.LecturerID != claimsUser.UserID {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, errors.New("unauthorized: lecturer ID mismatch"), pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.Delete(claims)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) GetAssignmentLecturer(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.GetAssignmentLecturer"

	claimsUser, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claimsUser == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetAssignmentLecturer(claimsUser.UserID)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) Submissions(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Submit"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(user.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	var req asgn.SubmissionRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.rctx.Logger.LogUserError(user.Email, err, pkgerr.ErrBodyParse.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	if time.Now().After(assignment.Deadline) {
		err := errors.New("The deadline has passed")
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrAssignmentDeadlinePassed.Message, pkgerr.ErrAssignmentDeadlinePassed.Details)
	}

	submission := models.Submission{
		AssignmentID: assignment.ID,
		StudentID:    user.UserID,
		FileURL:      req.FileURL,
	}

	respRepo, code, opRepo, err, msg, details := h.asgnRepo.Submissions(&submission, *assignment)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    respRepo,
	})
}

func (h *AssignmentsHandler) GetSubmission(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.GetSubmissions"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(user.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if user.UserID != assignment.ID {
		err := errors.New("The user does not own this assignment")
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetSubmission(assignment.ID)
	if err != nil {
		h.rctx.Logger.LogUserError(user.Email, err, msg)
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) GetSubmissions(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.GetSubmissions"

	user, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || user == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(user.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if user.UserID != assignment.ID {
		err := errors.New("The user does not own this assignment")
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetSubmissions(assignment.ID)
	if err != nil {
		h.rctx.Logger.LogUserError(user.Email, err, msg)
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) SubmissionGrade(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.SubmissionGrade"

	var req asgn.SubmissionGradeRequest

	submission, err := getSubmissionClaims(ctx)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingSubmissionID.Message, pkgerr.ErrMissingSubmissionID.Details)
	}

	lecturer, err := getUserClaims(ctx)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	if err := ctx.BodyParser(&req); err != nil {
		h.rctx.Logger.LogUserError(lecturer.Email, err, pkgerr.ErrBodyParse.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	req = asgn.SubmissionGradeRequest{
		SubmissionID: &submission.ID,
		LecturerID:   &lecturer.UserID,
		StatusID:     req.StatusID,
		Grade:        req.Grade,
		Notes:        req.Notes,
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.SubmissionGrade(req, submission)
	if err != nil {
		h.rctx.Logger.LogUserError(lecturer.Email, err, msg)
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func getUserClaims(ctx *fiber.Ctx) (
	*jwt.Claims,
	error,
) {
	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		return nil, errors.New("missing or invalid user claims")
	}

	return claims, nil
}

func getAssignmentClaims(ctx *fiber.Ctx) (
	*models.Assignment,
	error,
) {
	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		return nil, errors.New("missing or invalid assignment claims")
	}
	return assignment, nil
}

func getSubmissionClaims(ctx *fiber.Ctx) (
	*models.Submission,
	error,
) {
	sub, ok := ctx.Locals("submission").(*models.Submission)
	if !ok || sub == nil {
		return nil, errors.New("missing or invalid submission claims")
	}
	return sub, nil
}
