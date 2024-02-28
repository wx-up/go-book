package web

import (
	"time"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/wx-up/go-book/internal/global"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web/user"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	store := sessions.NewCookieStore([]byte("go-book"))
	engine.Use(sessions.Sessions("ssid", store))
	registerUserRoutes(engine)
}

func registerUserRoutes(engine *gin.Engine) {
	ug := engine.Group("/users")
	// 依赖注入的写法，遵循一个原则：我要用的东西我不会在内部自己初始化，由外部传入
	userDao := dao.NewUserDAO(global.DB)
	repo := repository.NewUserRepository(userDao)
	svc := service.NewUserService(repo)
	u := user.NewHandler(svc)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}
