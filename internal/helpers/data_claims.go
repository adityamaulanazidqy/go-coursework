package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
)

func GetUserClaims(ctx *fiber.Ctx) (
	*jwt.Claims,
	error,
) {
	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok || claims == nil {
		return nil, errors.New("missing or invalid user claims")
	}

	return claims, nil
}

func GetAssignmentClaims(ctx *fiber.Ctx) (
	*models.Assignment,
	error,
) {
	assignment, ok := ctx.Locals("assignment").(*models.Assignment)
	if !ok || assignment == nil {
		return nil, errors.New("missing or invalid assignment claims")
	}
	return assignment, nil
}

func GetSubmissionClaims(ctx *fiber.Ctx) (
	*models.Submission,
	error,
) {
	sub, ok := ctx.Locals("submission").(*models.Submission)
	if !ok || sub == nil {
		return nil, errors.New("missing or invalid submission claims")
	}
	return sub, nil
}
