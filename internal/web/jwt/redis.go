package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenKey  = []byte("go-book:access-token")
	RefreshTokenKey = []byte("go-book:refresh-token")
)

type RedisJwtHandler struct {
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
	accessTokenKey     []byte
	refreshTokenKey    []byte
	cmd                redis.Cmdable
}

func NewRedisJwtHandler(cmd redis.Cmdable) *RedisJwtHandler {
	return &RedisJwtHandler{
		accessTokenExpire:  time.Hour * 2,
		refreshTokenExpire: time.Hour * 24 * 7,
		accessTokenKey:     AccessTokenKey,
		refreshTokenKey:    RefreshTokenKey,
		cmd:                cmd,
	}
}

func (r *RedisJwtHandler) SetAccessToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.accessTokenExpire)), // 设置有效期
		},
		Uid:  uid,
		Ssid: ssid,
	})
	jwtToken, err := token.SignedString(r.accessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", jwtToken)
	return nil
}

func (r *RedisJwtHandler) SetRefreshToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	return r.setRefreshToken(ctx, uid, ssid)
}

func (r *RedisJwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.refreshTokenExpire)), // 设置有效期
		},
		Uid:       uid,
		Ssid:      ssid,
		IsRefresh: true,
	})
	jwtToken, err := token.SignedString(r.refreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", jwtToken)
	return nil
}

func (r *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	// 设置header 值为空，前端拿到空值之后会更新本地存储的 x-jwt-token 和 x-refresh-token，从而达到删除的效果
	// 这一步应该和前端协商好
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	val, _ := ctx.Get("claims")
	claims, ok := val.(*TokenClaim)
	if !ok {
		return errors.New("claims is not TokenClaim type")
	}

	// ssid 的有效期需要和长token一致
	return r.cmd.Set(ctx, claims.Ssid, "", r.refreshTokenExpire).Err()
}

func (r *RedisJwtHandler) ExtractToken(ctx *gin.Context) (string, error) {
	// 处理 jwt
	jwtOriStr := ctx.Request.Header.Get("Authorization")
	if jwtOriStr == "" {
		return "", errors.New("authorization header is missing")
	}
	jwtSlice := strings.SplitN(jwtOriStr, " ", 2)
	if len(jwtSlice) != 2 {
		return "", errors.New("authorization header is invalid")
	}
	return jwtSlice[1], nil
}

func (r *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	// 查看 ssid 是否在redis中
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("users:sid:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

// ParseToken claims 需要是指针类型
func (r *RedisJwtHandler) ParseToken(ctx *gin.Context, claims jwt.Claims, key []byte) error {
	jwtStr, err := r.ExtractToken(ctx)
	if err != nil {
		return err
	}
	token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("token 无效")
	}
	return nil
}
