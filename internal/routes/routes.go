package routes

import (
	"go-coursework/internal/models"
	"go-coursework/internal/routes/auth"
)

func SetupRoutes(ctx *models.RouterContext) {
	apiV1 := ctx.App.Group("/api/v1")
	auth.Setup(apiV1, ctx)
}
