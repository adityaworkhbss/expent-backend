package category

import (
    "github.com/gin-gonic/gin"
    "expent-backend/internal/category/handler"
    "expent-backend/internal/category/service"
    "expent-backend/internal/category/repository"
    "expent-backend/internal/infrastructure/prisma"
)

// RegisterRoutes registers category routes under the given router group.
func RegisterRoutes(rg *gin.RouterGroup, prismaClient *prisma.PrismaClient) {
    repo := repository.NewRepository(prismaClient)
    svc := service.NewService(repo)
    h := handler.NewHandler(svc)

    // Assuming Auth middleware sets "userId" in context.
    catGroup := rg.Group("/categories")
    catGroup.GET("", h.ListCategories)          // GET /categories
    catGroup.POST("", h.CreateCategory)        // POST /categories
    catGroup.DELETE("/:id", h.DeleteCategory) // DELETE /categories/:id
}
