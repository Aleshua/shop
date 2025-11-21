package router

import (
	l "shared/logger"

	"gateway/internal/config"
	locUtils "gateway/internal/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	logger l.ILogger
	router *gin.Engine
}

func New(config config.Services, logger l.ILogger) *GinRouter {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(requestid.New())

	r := &GinRouter{
		logger: logger,
		router: ginRouter,
	}

	r.initRoutes(config.AuthServiceUrl)

	return r
}

func (r GinRouter) Run(port string) {
	r.logger.With(l.NewField("port", port), l.NewField("mode", gin.Mode())).Info("сервер запускается")

	if err := r.router.Run(":" + port); err != nil {
		r.logger.With(l.NewField("error", err.Error())).Error("не удалось запустиь сервер")
	}
}

func (r GinRouter) initRoutes(authServiceUrl string) {
	authProxyHandler := locUtils.NewReverseProxyHandler(authServiceUrl, "/auth")
	authGroup := r.router.Group("/auth")
	{
		authGroup.GET("/ping", authProxyHandler)
		authGroup.POST("/register", authProxyHandler)
		authGroup.POST("/login", authProxyHandler)
		authGroup.POST("/logout", authProxyHandler)
		authGroup.POST("/email/confirm", authProxyHandler)
		authGroup.POST("/email/resend", authProxyHandler)
		authGroup.POST("/token/refresh", authProxyHandler)
	}
}
