package asgn

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go-coursework/internal/dto/asgn"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"go-coursework/pkg/jwt"
	"strings"
)

func (h *AssignmentsHandler) Comment(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.Comment"

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

	var req asgn.AssignmentCommentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	if strings.TrimSpace(req.Content) == "" {
		err := errors.New(pkgerr.ErrEmptyContent.Message)
		h.rctx.Logger.LogUserError(claimsUser.Email, err, pkgerr.ErrEmptyContent.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrEmptyContent.Message, pkgerr.ErrEmptyContent.Details)
	}

	req.UserID = claimsUser.UserID
	req.AssignmentID = claims.ID

	resp, code, opRepo, err, msg, details := h.asgnRepo.Comment(&req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) GetComments(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.GetComments"

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.GetCommentsByAssignmentID(assignment.ID)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}

func (h *AssignmentsHandler) DeleteComment(ctx *fiber.Ctx) error {
	const op = "handler.AssignmentsHandler.DeleteComment"

	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		err := errors.New(pkgerr.ErrMissingClaims.Message)
		h.rctx.Logger.LogUserError("-", err, pkgerr.ErrMissingClaims.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		err := errors.New(pkgerr.ErrAssignmentNotFound.Message)
		h.rctx.Logger.LogUserError(claims.Email, err, pkgerr.ErrAssignmentNotFound.Message)
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrAssignmentNotFound.Message, pkgerr.ErrAssignmentNotFound.Details)
	}

	if assignment.LecturerID != claims.UserID {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusUnauthorized, op, errors.New("unauthorized: lecturer ID mismatch"), pkgerr.ErrUnauthorized.Message, pkgerr.ErrUnauthorized.Details)
	}

	var req = asgn.DeleteComment{
		UserID:       claims.UserID,
		AssignmentID: assignment.ID,
		CommentID:    ctx.Params("comment_id"),
	}

	resp, code, opRepo, err, msg, details := h.asgnRepo.DeleteComment(&req)
	if err != nil {
		return h.rctx.Logger.LogRequestError(ctx, code, opRepo, err, msg, details)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data":    resp,
	})
}
