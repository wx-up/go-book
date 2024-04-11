package web

import (
	"fmt"
	"net/http"

	"github.com/wx-up/go-book/pkg/ginx"

	"github.com/wx-up/go-book/internal/web/jwt"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-gonic/gin"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (h *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/articles")
	g.POST("/save", h.Save)
	g.POST("/publish", ginx.WrapHandleWithClaim[PublishArticleReq, jwt.UserClaim]("claims", h.Publish))
	g.POST("/withdraw", h.Withdraw)
}

type PublishArticleReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Withdraw 撤回，仅自己可见
func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
}

// Publish 发布
func (h *ArticleHandler) Publish(ctx *gin.Context, req PublishArticleReq, claim jwt.UserClaim) (Result, error) {
	id, err := h.svc.Publish(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claim.Uid,
		},
	})
	if err != nil {
		// 包装一下错误，日志统一去 wrap 中打印
		return Result{
			Code: 5,
			Msg:  "服务器错误，请稍后再试",
		}, fmt.Errorf("发布文章错误：%w", err)
	}
	return Result{
		Code: 0,
		Msg:  "发布成功",
		Data: map[string]any{
			"id": id,
		},
	}, nil
}

// Save 新增或者编辑
func (h *ArticleHandler) Save(c *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Code: -1,
			Msg:  "参数错误",
		})
		return
	}
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, Result{
			Code: -1,
			Msg:  "参数错误",
		})
		return
	}
	claim := c.Value("claims").(jwt.UserClaim)
	// 调用 articleService
	id, err := h.svc.Save(c, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claim.Uid,
		},
	})
	if err == service.ErrArticleNotFound {
		c.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "参数错误",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "服务器错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "保存成功",
		Data: map[string]any{
			"id": id,
		},
	})
}