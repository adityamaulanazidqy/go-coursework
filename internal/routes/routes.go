package routes

import (
	"go-coursework/internal/models"
	"go-coursework/internal/routes/asgn"
	"go-coursework/internal/routes/auth"
	"go-coursework/internal/routes/otp"
	"go-coursework/internal/routes/settings"
)

func SetupRoutes(rctx *models.RouterContext) {
	apiV1 := rctx.App.Group("/api/v1")
	auth.Setup(apiV1, rctx)
	settings.Setup(apiV1, rctx)
	otp.Setup(apiV1, rctx)
	asgn.Setup(apiV1, rctx)
}
