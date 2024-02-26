package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web/user"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	registerUserRoutes(engine)
}

func registerUserRoutes(engine *gin.Engine) {
	ug := engine.Group("/users")
	u := user.NewHandler()
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}
