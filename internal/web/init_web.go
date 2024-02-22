package web

import (
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web/user"
)

func RegisterRoutes(engine *gin.Engine) {
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
