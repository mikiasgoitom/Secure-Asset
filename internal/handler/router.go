package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/handler/middleware"
)

type Router struct{
	userhandler *UserHandler
	AssetHandler *AssetHandler
	jwtService contract.IJWTService
}

func NewRouter(userUsecase contract.IUserUsecase, assetUsecase contract.IAssetUsecase, jwtService contract.IJWTService, logger contract.ILogger) *Router {
	return &Router{
		userhandler: NewUserHandler(userUsecase, logger),
		AssetHandler: NewAssetHandler(assetUsecase, logger),
		jwtService: jwtService,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	user := v1.Group("/user") 
	{
		user.POST("/register", r.userhandler.RegisterUser)
		user.POST("/login", r.userhandler.LoginUser)
	}
	asset := v1.Group("/asset")
	asset.Use(middleware.AuthMiddleware(r.jwtService))
	{
		asset.POST("/create", r.AssetHandler.CreateAsset)
	}
}