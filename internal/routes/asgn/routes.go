package asgn

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/asgn"
	"go-coursework/internal/models"
	asgnmiddleware "go-coursework/pkg/asgn"
	"go-coursework/pkg/jwt"
)

func Setup(router fiber.Router, rctx *models.RouterContext) {
	controller := asgn.NewAssignmentsHandler(rctx)

	assignmentGroup := router.Group("/assignments")
	{
		assignmentGroup.Get("/lecturer", jwt.Middleware("Lecturer"), controller.GetAssignmentLecturer)

		assignmentGroup.Post("", jwt.Middleware("Lecturer", "Admin"), controller.Post)
		assignmentGroup.Get("/all", jwt.Middleware("Student", "Admin"), controller.GetAll)
		assignmentGroup.Get("/:id", jwt.Middleware("Lecturer", "Admin", "Student"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Get)
		assignmentGroup.Put("/:id", jwt.Middleware("Lecturer", "Admin"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Update)
		assignmentGroup.Delete("/:id", jwt.Middleware("Lecturer", "Admin"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Delete)

		assignmentGroup.Post("/:id/submissions", jwt.Middleware("Student"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Submissions)

		commentGroup := assignmentGroup.Group("/:id/comments")
		{
			commentGroup.Post("", jwt.Middleware("Student", "Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Comment)
			commentGroup.Get("", jwt.Middleware("Student", "Lecturer", "Admin"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.GetComments)
			commentGroup.Delete("/:comment_id", jwt.Middleware("Lecturer", "Admin"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.DeleteComment)
		}
	}
}
