package handler

import (
	"jokegen/internal/service"

	_ "jokegen/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Handler struct {
	Users Users
}

func NewHandler(user *service.UsersService) *Handler {
	return &Handler{Users: user}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/user")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/check-user", h.checkUser)
	}
	jokes := router.Group("/jokes", h.userIdentity)
	{
		jokes.GET("/", h.GetJoke)
	}
	return router
}
