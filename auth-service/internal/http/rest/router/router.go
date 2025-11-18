package router

import (
	c "auth/internal/config"
	l "auth/internal/logger"
	ucPorts "auth/internal/usecases/ports"

	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	port string

	logger l.ILogger
	router *gin.Engine
	uc     ucPorts.IAuthUseCases
}

func New(params c.HTTP, logger l.ILogger, uc ucPorts.IAuthUseCases) *GinRouter {
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())

	r := GinRouter{
		port:   params.Port,
		logger: logger,
		router: ginRouter,
		uc:     uc,
	}

	r.initRoutes()

	return &r
}

func (r *GinRouter) Run() {
	r.router.Run(":" + r.port)
}

func (r *GinRouter) GetGinRouter() *gin.Engine {
	return r.router
}

func (r *GinRouter) initRoutes() {
	r.router.GET("ping", r.ping)
	r.router.POST("register", r.register)
	r.router.POST("login", r.login)
	r.router.POST("logout", r.logout)
	r.router.POST("email/confirm", r.confirmEmail)
	r.router.POST("email/resend", r.resendCode)
	r.router.POST("token/refresh", r.refreshToken)
}
