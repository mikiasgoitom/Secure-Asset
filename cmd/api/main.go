package main

import (
	"context"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/handler"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/config"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/logger"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/email"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/repository"
	"github.com/mikiasgoitom/Secure-Asset/internal/infrastructure/security"
	"github.com/mikiasgoitom/Secure-Asset/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// --- Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// log.Println(cfg, "sdfg", cfg.DatabaseName)
	// --- Database Connection ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(cfg.DatabaseURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer client.Disconnect(ctx)

	db := client.Database(cfg.DatabaseName)
	log.Println("Successfully connected to MongoDB.")

	// casbin initialization
	log.Println("database", cfg.DatabaseName)
	casbinAdapter, err := mongodbadapter.NewAdapterWithCollectionName(clientOptions, cfg.DatabaseName, cfg.CasbinCollection)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin adapter: %v", err)
	}

	casbinEnforcer, err := casbin.NewEnforcer(cfg.CasbinModelPath, casbinAdapter)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}
	if err := casbinEnforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load Casbin policy: %v", err)
	}
	log.Println("Casbin enforcer initialized successfully.")

	// --- Dependency Injection (Wiring the application together) ---

	// Infrastructure Layer
	zapLogger, err := logger.NewZapAdapter(cfg.Production)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	userRepo := repository.NewUserRepository(db, cfg.UserCollection)
	assetRepo := repository.NewAssetRepository(db, cfg.AssetCollection)
	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer)
	emailService := email.NewMailtrapService(cfg.MailtrapHost, cfg.MailtrapPort, cfg.MailtrapUser, cfg.MailtrapPass, cfg.MailFrom)

	// Usecase Layer
	userUsecase := usecase.NewUserUsecase(userRepo, zapLogger, jwtService, casbinEnforcer, emailService)
	assetUsecase := usecase.NewAssetUsecase(assetRepo, userRepo, zapLogger)
	// handler layer
	assetHandler := handler.NewAssetHandler(assetUsecase, zapLogger)
	userHandler := handler.NewUserHandler(userUsecase, zapLogger)
	// Presentation Layer (Handler & Router)
	appRouter := handler.NewRouter(userHandler, assetHandler, jwtService, zapLogger, casbinEnforcer)

	// --- Initialize Gin Router ---
	ginEngine := gin.Default()

	// Global middleware should be registered before routes so all routes get it.
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true // Development only; tighten for production
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	ginEngine.Use(cors.New(corsCfg))

	appRouter.SetupRoutes(ginEngine)
	// --- Start Server ---
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
