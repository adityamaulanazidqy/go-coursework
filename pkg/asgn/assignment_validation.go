package asgn

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"gorm.io/gorm"
)

func AssignmentExistMiddleware(rctx *models.RouterContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		asgnID := ctx.Params("id")
		if asgnID == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(pkgerr.ErrMissingAsgnID)
		}

		var assignment models.Assignment
		if err := rctx.DB.First(&assignment, "id = ?", asgnID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(pkgerr.ErrAssignmentNotFound)
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(pkgerr.ErrInternalServer)
		}

		ctx.Locals("assignment", &assignment)
		return ctx.Next()
	}
}

func SubmissionExistMiddleware(rctx *models.RouterContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		submissionID := ctx.Params("submission_id")
		if submissionID == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(pkgerr.ErrMissingSubmissionID)
		}

		var submission models.Submission
		if err := rctx.DB.First(&submission, "id = ?", submissionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(pkgerr.ErrSubmissionStatusNotFound)
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(pkgerr.ErrInternalServer)
		}

		ctx.Locals("submission", &submission)
		return ctx.Next()
	}
}
