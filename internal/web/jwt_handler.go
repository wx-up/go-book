package web

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wx-up/go-book/internal/domain"
)

type jwtHandler struct{}

type TokenClaim struct {
	jwt.RegisteredClaims
	Uid       int64
	IsRefresh bool
	Ssid      string
}

type UserClaim = TokenClaim

func (h *jwtHandler) setLoginToken(ctx *gin.Context, u domain.User) error {
	ssid := uuid.New().String()
	err := h.setJwtToken(ctx, u, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, u, ssid)
	if err != nil {
		return err
	}
	return nil
}

// setRefreshToken 设置 refresh_token
func (h *jwtHandler) setRefreshToken(ctx *gin.Context, u domain.User, ssid string) error {
	// jwt 保持登陆状态
	jwtToken, err := h.generateRefreshToken(u, ssid)
	if err != nil {
		return errors.New("系统错误")
	}
	ctx.Header("x-refresh-token", jwtToken)
	return nil
}

func (h *jwtHandler) generateRefreshToken(u domain.User, ssid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "go-book",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)), // 设置有效期
		},
		Uid:       u.Id,
		IsRefresh: true,
		Ssid:      ssid,
	})
	jwtToken, err := token.SignedString([]byte("go-book"))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

// setJwtToken 设置 access_token
func (h *jwtHandler) setJwtToken(ctx *gin.Context, u domain.User, ssid string) error {
	// jwt 保持登陆状态
	jwtToken, err := h.generateJwtToken(u, ssid)
	if err != nil {
		return errors.New("系统错误")
	}
	ctx.Header("x-jwt-token", jwtToken)
	return nil
}

func (h *jwtHandler) generateJwtToken(u domain.User, ssid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "go-book",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)), // 设置有效期
		},
		Uid:  u.Id,
		Ssid: ssid,
	})
	jwtToken, err := token.SignedString([]byte("go-book"))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
