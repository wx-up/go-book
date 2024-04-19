package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/wx-up/go-book/pkg/slice"

	"github.com/wx-up/go-book/pkg/ginx"

	"github.com/wx-up/go-book/internal/web/jwt"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-gonic/gin"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	incrSvc service.InteractiveService
	l       logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger, incrSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		incrSvc: incrSvc,
	}
}

func (h *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/articles")
	g.POST("/save", h.Save)
	g.POST("/publish", ginx.WrapHandleWithReqAndClaim[PublishArticleReq, jwt.UserClaim]("claims", h.Publish))
	g.POST("/withdraw", h.Withdraw)
	g.GET("/", ginx.WrapHandleWithReqAndClaim[ArticleListReq, jwt.UserClaim]("claims", h.List))
	g.GET("/detail/:id", ginx.WrapHandleWithClaim[jwt.UserClaim]("claims", h.Detail))

	// 线上库，也就是用户端访问
	gp := g.Group("/pub")
	gp.GET("/:id",
		ginx.WrapHandleWithClaim[jwt.UserClaim]("claims", h.PublishedDetail),
	)
}

func (h *ArticleHandler) PublishedDetail(ctx *gin.Context, claims jwt.UserClaim) (Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Result{
			Code: -1,
			Msg:  "参数错误",
		}, err
	}

	// 因为获取文章主体和指标信息并不是互相依赖的，所以可以并发处理
	var eg errgroup.Group
	var detail domain.Article
	eg.Go(func() error {
		detail, err = h.svc.PublishedDetail(ctx, id)
		return err
	})
	var inter domain.Interactive
	eg.Go(func() error {
		inter, err = h.incrSvc.Get(ctx, "articles", detail.Id, claims.Uid)
		// 容错
		if err != nil {
			h.l.Error("获取文章指标信息失败", logger.Error(err), logger.Int64("aid", detail.Id))
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		return Result{
			Code: 5,
			Msg:  "服务器错误",
		}, err
	}

	// 增加阅读计数
	//go func() {
	//	er := h.incrSvc.IncrReadCount(ctx, "articles", id)
	//	if er != nil {
	//		h.l.Error("增加阅读计数失败", logger.Error(er), logger.Int64("aid", id))
	//	}
	//}()

	return Result{
		Data: ArticleVO{
			Id:         detail.Id,
			Title:      detail.Title,
			Abstract:   detail.Abstract(),
			Content:    detail.Content,
			AuthorId:   detail.Author.Id,
			AuthorName: detail.Author.Name,
			CreateTime: detail.CreateTime.Format(time.DateTime),
			UpdateTime: detail.UpdateTime.Format(time.DateTime),
			LikeCnt:    inter.LikeCnt,
			CollectCnt: inter.CollectCnt,
			ReadCnt:    inter.ReadCnt,
			Liked:      inter.Liked,
			Collected:  inter.Collected,
			Status:     detail.Status.ToUint8(),
		},
	}, nil
}

type PublishArticleReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *ArticleHandler) Detail(ctx *gin.Context, claim jwt.UserClaim) (Result, error) {
	// 路由参数只能通过如下的方式获取和校验，无法通过 wrap 来优化
	// 可以和前端约定，全部用 post 来请求，那么所有接口都可以使用 wrap 优化
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Result{
			Code: -1,
			Msg:  "参数错误",
		}, err
	}
	obj, err := h.svc.Detail(ctx, id)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "服务器错误",
		}, err
	}
	if obj.Author.Id != claim.Uid {
		// 非法访问文章，这时候要上报这种非法用户
		return Result{
			Code: 4,
			// 不需要告诉前端发生了什么
			Msg: "输入有误",
		}, fmt.Errorf("非法访问文章，用户ID：%d，文章作者ID：%d", claim.Uid, obj.Author.Id)
	}
	return Result{
		Code: 0,
		Msg:  "获取成功",
		Data: ArticleVO{
			Id:       obj.Id,
			Title:    obj.Title,
			Abstract: obj.Abstract(),
			Content:  obj.Content,
			// AuthorId:   val.Author.Id,
			// AuthorName: "",
			CreateTime: obj.CreateTime.Format(time.DateTime),
			UpdateTime: obj.CreateTime.Format(time.DateTime),
			Status:     obj.Status.ToUint8(),
		},
	}, nil
}

// List 作者文章列表
func (h *ArticleHandler) List(ctx *gin.Context, req ArticleListReq, claims jwt.UserClaim) (Result, error) {
	objs, err := h.svc.List(ctx, claims.Uid, req.Page, req.Size)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "服务器错误",
		}, fmt.Errorf("获取文章列表错误：%w", err)
	}
	return Result{
		Code: 0,
		Msg:  "获取成功",
		Data: slice.Map[domain.Article, ArticleVO](objs, func(idx int, val domain.Article) ArticleVO {
			return ArticleVO{
				Id:       val.Id,
				Title:    val.Title,
				Abstract: val.Abstract(),
				// Content:    val.Content,
				// AuthorId:   val.Author.Id,
				// AuthorName: "",
				CreateTime: val.CreateTime.Format(time.DateTime),
				UpdateTime: val.CreateTime.Format(time.DateTime),
				Status:     val.Status.ToUint8(),
			}
		}),
	}, nil
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
