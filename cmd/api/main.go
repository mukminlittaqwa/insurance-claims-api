package main

import (
    "context"
    "insurance-claims-api/internal/config"
    "insurance-claims-api/internal/handlers"
    "insurance-claims-api/internal/middleware"
    "insurance-claims-api/internal/repositories"
    "insurance-claims-api/internal/services"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gin-contrib/cors"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client     *mongo.Client
    claimService services.ClaimService
    authService  services.AuthService
)

func main() {
    config.LoadConfig()

    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.AppConfig.MongoURI))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err = client.Ping(ctx, nil); err != nil {
        log.Fatal("Cannot connect to MongoDB:", err)
    }
    log.Println("Connected to MongoDB Atlas!")

    userRepo := repositories.NewUserRepository(client)
    claimRepo := repositories.NewClaimRepository(client)

    authService = services.NewAuthService(userRepo)
    claimService = services.NewClaimService(claimRepo, userRepo)

    r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000","http://localhost:3001",},
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.POST("/api/v1/login", handlers.Login(authService))

    authRoutes := r.Group("/api/v1")
    authRoutes.Use(middleware.AuthMiddleware())

    authRoutes.POST("/claims", handlers.CreateClaim(claimService))
    authRoutes.GET("/claims", handlers.GetMyClaims(claimService))
    authRoutes.GET("/claims/all", handlers.GetAllClaims(claimService)) // hanya verifier & approver
    authRoutes.GET("/claims/:id", handlers.GetClaimByID(claimService))
    authRoutes.PATCH("/claims/:id", handlers.UpdateClaim(claimService))
    authRoutes.DELETE("/claims/:id", handlers.DeleteClaim(claimService))
    authRoutes.PATCH("/claims/:id/submit", handlers.SubmitClaim(claimService))
    authRoutes.PATCH("/claims/:id/review", handlers.ReviewClaim(claimService))
    authRoutes.PATCH("/claims/:id/approve", handlers.ApproveClaim(claimService))
    authRoutes.PATCH("/claims/:id/reject", handlers.RejectClaim(claimService))

    r.Run(":" + config.AppConfig.Port)
}