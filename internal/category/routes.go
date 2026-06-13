package category

import (
	"expent-backend/internal/category/handler"
	"expent-backend/internal/category/repository"
	"expent-backend/internal/category/service"
	"expent-backend/internal/infrastructure/prisma"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers category routes under the given router group.
func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
	repo := repository.NewRepository(prismaClient)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	// Assuming Auth middleware sets "userId" in context.
	catGroup := rg.Group("/categories")
	catGroup.GET("", h.ListCategories)        // GET /categories
	catGroup.POST("", h.CreateCategory)       // POST /categories
	catGroup.PUT("/:id", h.UpdateCategory)    // PUT /categories/:id
	catGroup.DELETE("/:id", h.DeleteCategory) // DELETE /categories/:id
}
