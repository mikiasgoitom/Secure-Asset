package handler

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/handler/middleware"
)

type Router struct{
	userhandler *UserHandler
	AssetHandler *AssetHandler
	jwtService contract.IJWTService
	casbinEnforcer *casbin.Enforcer
}

func NewRouter(userHandler *UserHandler, assetHandler *AssetHandler, jwtService contract.IJWTService, logger contract.ILogger, enforcer *casbin.Enforcer) *Router {
	return &Router{
		userhandler: userHandler,
		AssetHandler: assetHandler,
		jwtService: jwtService,
		casbinEnforcer: enforcer,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	user := v1.Group("/user") 
	// user.Use(middleware.Authentication(r.jwtService))
	{
		user.POST("/register", r.userhandler.RegisterUser)
		user.POST("/login", r.userhandler.LoginUser)
	}
	asset := v1.Group("/asset")
	asset.Use(middleware.Authentication(r.jwtService))
	asset.Use(middleware.Authorization(r.casbinEnforcer))
	{
		asset.POST("/create", r.AssetHandler.CreateAsset)
		asset.GET("/:id", r.AssetHandler.GetAsset)
	}
}