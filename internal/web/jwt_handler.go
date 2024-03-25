package web

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wx-up/go-book/internal/domain"
)

type jwtHandler struct{}

type TokenClaim struct {
	jwt.RegisteredClaims
	Uid       int64
	IsRefresh bool
}

type UserClaim = TokenClaim

// setRefreshToken 设置 refresh_token
func (h *jwtHandler) setRefreshToken(ctx *gin.Context, u domain.User) error {
	// jwt 保持登陆状态
	jwtToken, err := h.generateRefreshToken(u)
	if err != nil {
		return errors.New("系统错误")
	}
	ctx.Header("x-refresh-token", jwtToken)
	return nil
}

func (h *jwtHandler) generateRefreshToken(u domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "go-book",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)), // 设置有效期
		},
		Uid:       u.Id,
		IsRefresh: true,
	})
	jwtToken, err := token.SignedString([]byte("go-book"))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

// setJwtToken 设置 access_token
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
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
