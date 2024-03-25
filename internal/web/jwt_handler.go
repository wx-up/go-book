package web

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wx-up/go-book/internal/domain"
	"github.com/wx-up/go-book/internal/web/middleware"
)

type jwtHandler struct{}

func (h *jwtHandler) setJwtToken(ctx *gin.Context, u domain.User) error {
	// jwt 保持登陆状态
	jwtToken, err := h.generateJwtToken(u)
	if err != nil {
		return errors.New("系统错误")
	}
	ctx.Header("x-jwt-token", jwtToken)
	return nil
}

func (h *jwtHandler) generateJwtToken(u domain.User) (string, error) {
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
