package user

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewHandler() *Handler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &Handler{
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (h *Handler) SignUp(ctx *gin.Context) {
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

	// 注册成功
	ctx.String(http.StatusOK, "注册成功")
}

func (h *Handler) Login(ctx *gin.Context) {
}

func (h *Handler) Edit(ctx *gin.Context) {
}

func (h *Handler) Profile(ctx *gin.Context) {
}
