package web

import (
	"time"

	"github.com/wx-up/go-book/internal/web/middleware"

	"github.com/wx-up/go-book/internal/global"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web/user"
)

func RegisterRoutes(engine *gin.Engine) {
	// 跨域
	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"X-Jwt-Token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// session 插件
	//store, err := redis.NewStore(
	//	16,
	//	"tcp",
	//	"localhost:7379",
	//	"",
	//	[]byte("Kv5mvUKCUDmGRC2XRZI622fWvazQaHCB"),
	//	[]byte("bOCdz7AdaFiRTF8kiLVxY7I8BHn49dPh"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//engine.Use(sessions.Sessions("ssid", store))

	// 登陆插件
	engine.Use(middleware.NewLoginJwtMiddlewareBuilder().Build())

	// 注册业务路由
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
