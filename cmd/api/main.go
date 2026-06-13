package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "go.uber.org/zap"

    "expent-backend/configs"
    "expent-backend/internal/infrastructure/prisma"
    "expent-backend/internal/middleware"
    "expent-backend/internal/auth"
    "expent-backend/internal/account"
    "expent-backend/internal/budget"
    "expent-backend/internal/category"
    "expent-backend/internal/dashboard"
    "expent-backend/internal/transaction"
    parsetransaction "expent-backend/internal/parse_transaction"
    usercustomization "expent-backend/internal/user_customization"
    "expent-backend/internal/emi"
)

func main() {
    // Load environment variables first
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found: %v", err)
    }
    // Load configuration into global variable
    configs.LoadConfig()

    // Initialize logger (zap)
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()
    zap.ReplaceGlobals(logger)

    // Initialize Prisma client
    prismaClient, err := prisma.NewClient(configs.AppConfig.DATABASE_URL)
    if err != nil {
        logger.Fatal("Failed to initialize Prisma client", zap.Error(err))
    }
    defer prismaClient.Prisma.Disconnect()

    // Create Gin router
    r := gin.New()
    // Global middleware
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())

    // Public routes
    api := r.Group("/" + configs.AppConfig.API_PREFIX)
    auth.RegisterRoutes(api, prismaClient)

    // Protected routes (requires JWT)
    protected := api.Group("")
    protected.Use(middleware.Auth())
    
    account.RegisterRoutes(protected, prismaClient)
    budget.RegisterRoutes(protected, prismaClient)
    category.RegisterRoutes(protected, prismaClient)
    transaction.RegisterRoutes(protected, prismaClient)
    parsetransaction.RegisterRoutes(protected)
    dashboard.RegisterRoutes(protected)
    usercustomization.RegisterRoutes(protected, prismaClient)
    emi.RegisterRoutes(protected, prismaClient)

    // Start server
    port := configs.AppConfig.PORT
    if port == "" {
        port = "3000"
    }
    logger.Info("Starting server", zap.String("port", port))
    if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
        logger.Fatal("Server stopped unexpectedly", zap.Error(err))
    }
}
