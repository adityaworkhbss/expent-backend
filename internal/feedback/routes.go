package feedback

import (
	"expent-backend/internal/feedback/handler"
	"expent-backend/internal/feedback/repository"
	"expent-backend/internal/feedback/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	feedbackGroup := r.Group("/feedback")
	feedbackGroup.POST("", h.SubmitFeedback)
}
