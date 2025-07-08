package asgn

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

		assignmentGroup.Post("", jwt.Middleware("Lecturer"), controller.Post)
		assignmentGroup.Get("/all", jwt.Middleware("Student"), controller.GetAll)
		assignmentGroup.Get("/:id", jwt.Middleware("Lecturer", "Student"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Get)
		assignmentGroup.Put("/:id", jwt.Middleware("Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Update)
		assignmentGroup.Delete("/:id", jwt.Middleware("Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Delete)

		assignmentGroup.Post("/:id/submissions", jwt.Middleware("Student"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Submissions)
		assignmentGroup.Put("/:id/submissions", jwt.Middleware("Student"), asgnmiddleware.SubmissionExistMiddleware(rctx), controller.UpdateSubmission)
		assignmentGroup.Get("/:id/submissions", jwt.Middleware("Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.GetSubmissions)
		assignmentGroup.Get("/:id/submission", jwt.Middleware("Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.GetSubmission)
		assignmentGroup.Post("/:submission_id/submissions/grade", jwt.Middleware("Lecturer"), asgnmiddleware.SubmissionExistMiddleware(rctx), controller.SubmissionGrade)

		commentGroup := assignmentGroup.Group("/:id/comments")
		{
			commentGroup.Post("", jwt.Middleware("Student", "Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.Comment)
			commentGroup.Get("", jwt.Middleware("Student", "Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.GetComments)
			commentGroup.Delete("/:comment_id", jwt.Middleware("Lecturer"), asgnmiddleware.AssignmentExistMiddleware(rctx), controller.DeleteComment)
		}

		searchGroup := assignmentGroup.Group("/search")
		{
			searchGroup.Get("/:id/submissions/ws",
				jwt.MiddlewareSocket("Lecturer"),
				asgnmiddleware.AssignmentExistMiddleware(rctx),
				websocket.New(controller.SearchSubmissionWS),
			)
		}
	}
}
