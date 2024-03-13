package web

import (
	"net/http"
	"time"

	"github.com/wx-up/go-book/internal/service/code"

	"github.com/wx-up/go-book/internal/web/middleware"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

const biz = "login"

type UserHandler struct {
	svc         *service.UserService
	codeSvc     *code.SmsCodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

var _ handler = (*UserHandler)(nil)

func NewUserHandler(svc *service.UserService, codeSvc *code.SmsCodeService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		codeSvc:     codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	ug := engine.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.POST("/profile", h.Profile)

	// 验证码发送
	ug.POST("/code/send", h.SendCode)
	// 验证码验证+登陆
	ug.POST("/code/verify", h.VerifyCode)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 邮箱验证
	ok, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	// 密码验证
	ok, err = h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式不正确")
		return
	}

	// 注册
	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "该邮箱已注册")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 注册成功
	ctx.String(http.StatusOK, "注册成功")
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 邮箱验证
	u, err := h.svc.Login(ctx, domain.User{Email: req.Email, Password: req.Password})
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 保持登陆状态
	//sess := sessions.Default(ctx)
	//sess.Set("uid", u.Id)
	//sess.Options(sessions.Options{
	//	MaxAge: 30 * 60, // 三十分钟
	//})
	//if err = sess.Save(); err != nil {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	// jwt 保持登陆状态
	jwtToken, err := h.generateJwtToken(u)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", jwtToken)

	ctx.String(http.StatusOK, "登陆成功")
}

func (h *UserHandler) generateJwtToken(u domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "go-book",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)), // 设置有效期
		},
		Uid: u.Id,
	})
	jwtToken, err := token.SignedString([]byte("go-book"))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (h *UserHandler) Edit(ctx *gin.Context) {
}

func (h *UserHandler) Profile(ctx *gin.Context) {
}

func (h *UserHandler) SendCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 需要正则表达式强验证
	if len(req.Phone) != 11 {
		ctx.JSON(http.StatusOK, Result{Msg: "手机号格式错误", Code: 4})
		return
	}

	// 发送验证码
	err := h.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "验证码发送成功"})
	case code.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{Msg: "验证码发送过于频繁，请稍后再试", Code: 2})
	default:
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误", Code: 5})
	}
}

func (h *UserHandler) VerifyCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err == code.ErrCodeVerifyTooMany {
		ctx.JSON(http.StatusOK, Result{Msg: "验证过于频繁，请稍后再试", Code: 5})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "验证码错误", Code: 5})
		return
	}

	// 登陆用户
	u, err := h.svc.FindOrCreateByPhone(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "验证码错误", Code: 5})
		return
	}
	jwtToken, err := h.generateJwtToken(u)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", jwtToken)
	ctx.String(http.StatusOK, "登陆成功")
}
