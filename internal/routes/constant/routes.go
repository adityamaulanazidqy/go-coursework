package constant

import (
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/handlers/constant"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
)

func Setup(router fiber.Router, rctx *models.RouterContext) {
	controller := constant.NewConstHandler(rctx)

	dataGroup := router.Group("data")
	semesterGroup := dataGroup.Group("semester")
	{
		semesterGroup.Get("", jwt.Middleware("Lecturer", "Admin"), controller.GetSemesters)
		semesterGroup.Post("", jwt.Middleware("Lecturer", "Admin"), controller.PostSemesters)
		semesterGroup.Delete("", jwt.Middleware("Lecturer", "Admin"), controller.DeleteSemester)
	}

	studyProgramGroup := dataGroup.Group("study-program")
	{
		studyProgramGroup.Get("", jwt.Middleware("Lecturer", "Admin"), controller.GetStudyPrograms)
		studyProgramGroup.Post("", jwt.Middleware("Lecturer", "Admin"), controller.PostStudyPrograms)
		studyProgramGroup.Delete("", jwt.Middleware("Lecturer", "Admin"), controller.DeleteProgram)
	}
}
