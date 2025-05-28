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

	req := asgn.AssignmentRequest{
		LecturerID:    claims.UserID,
		Title:         ctx.FormValue("title"),
		Description:   ctx.FormValue("description"),
		FileHeader:    fileHeader,
		MultipartFile: multipartFile,
		Deadline:      deadline,
		IsActive:      true,
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

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetAll()
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

	claims, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(claimsUser.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if claims.LecturerID != claimsUser.UserID {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, errors.New("unauthorized: lecturer ID mismatch"), pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	fileHeader, err := ctx.FormFile("file")
	var multipartFile multipart.File

	if err == nil && fileHeader != nil {
		multipartFile, err = fileHeader.Open()
		if err != nil {
			return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrFileOpen.Message, pkgerr.ErrFileOpen.Details)
		}
	} else {
		fileHeader = nil
		multipartFile = nil
	}

	deadlineStr := ctx.FormValue("deadline")
	var deadline time.Time

	if deadlineStr != "" {
		deadline, err = time.Parse("2006-01-02 15:04:05", deadlineStr)
		if err != nil {
			return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, "Invalid deadline format", nil)
		}
	} else {
		deadline = claims.Deadline
	}

	title := ctx.FormValue("title")
	if title != "" {
		claims.Title = title
	}

	description := ctx.FormValue("description")
	if description != "" {
		claims.Description = description
	}

	var req = asgn.AssignmentRequest{
		Title:         claims.Title,
		Description:   claims.Description,
		FileHeader:    fileHeader,
		MultipartFile: multipartFile,
		Deadline:      deadline,
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.Update(claims, &req)
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
