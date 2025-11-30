package api

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/handler"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/config"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/logger"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/repository"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/security"
	"github.com/mikiasgoitom/Secure-Asset/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// --- 1. Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- 2. Database Connection ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer client.Disconnect(ctx)

	db := client.Database(cfg.DatabaseName)
	log.Println("Successfully connected to MongoDB.")

	// --- 3. Dependency Injection (Wiring the application together) ---

	// Infrastructure Layer
	zapLogger, err := logger.NewZapAdapter(cfg.Production)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	userRepo := repository.NewUserRepository(db, cfg.UserCollection)
	assetRepo := repository.NewAssetRepository(db, cfg.AssetCollection)
	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer)

	// Usecase Layer
	userUsecase := usecase.NewUserUsecase(userRepo, zapLogger, jwtService)
	assetUsecase := usecase.NewAssetUsecase(assetRepo, userRepo, zapLogger)

	// Presentation Layer (Handler & Router)
	appRouter := handler.NewRouter(userUsecase, assetUsecase, jwtService, zapLogger)

	// --- 4. Initialize Gin Router ---
	ginEngine := gin.Default()
	appRouter.SetupRoutes(ginEngine)

	// --- 5. Start Server ---
	log.Println("Starting server on :" + cfg.ServerPort)
	if err := ginEngine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Application will gracefully shutdown and clean up resources here if needed.
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()
}
