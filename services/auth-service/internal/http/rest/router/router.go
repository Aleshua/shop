package router

import (
	l "shared/logger"

	mds "auth/internal/http/rest/middlewares"
	uc "auth/internal/usecases"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	logger l.ILogger
	router *gin.Engine
	uc     uc.IAuthUseCases
}

func New(logger l.ILogger, uc uc.IAuthUseCases) *GinRouter {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(requestid.New())
	ginRouter.Use(mds.LoggerMiddleware(logger))

	r := GinRouter{
		logger: logger,
		router: ginRouter,
		uc:     uc,
	}

	r.initRoutes()

	return &r
}

func (r *GinRouter) Run(port string) {
	r.logger.With(l.NewField("port", port), l.NewField("mode", gin.Mode())).Info("сервер запускается")

	if err := r.router.Run(":" + port); err != nil {
		r.logger.With(l.NewField("error", err.Error())).Error("не удалось запустиь сервер")
	}
}

func (r *GinRouter) initRoutes() {
	r.router.GET("/ping", r.ping)
	r.router.POST("/register", r.register)
	r.router.POST("/login", r.login)
	r.router.POST("/logout", r.logout)
	r.router.POST("/email/confirm", r.confirmEmail)
	r.router.POST("/email/resend", r.resendCode)
	r.router.POST("/token/refresh", r.refreshToken)
}
