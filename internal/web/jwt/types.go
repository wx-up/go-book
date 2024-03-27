package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetAccessToken(ctx *gin.Context, uid int64) error
	SetRefreshToken(ctx *gin.Context, uid int64) error
	ClearToken(ctx *gin.Context) error
	ExtractToken(ctx *gin.Context) (string, error)
	CheckSession(ctx *gin.Context, ssid string) error
	ParseToken(ctx *gin.Context, claims jwt.Claims, key []byte) error
}

type TokenClaim struct {
	jwt.RegisteredClaims
	Uid       int64
	IsRefresh bool
	Ssid      string
}

type UserClaim = TokenClaim
