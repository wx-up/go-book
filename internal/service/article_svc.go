package service

import (
	"context"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

var ErrArticleNotFound = repository.ErrArticleNotFound

//go:generate mockgen -destination=./mocks/article_svc.mock.go -package=svcmocks -source=article_svc.go
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	List(ctx context.Context, uid int64, page, size int64) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)

	PublishedDetail(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func (a *articleService) PublishedDetail(ctx context.Context, id int64) (domain.Article, error) {
	// 在这里调用 userService 进行 组装
	// TODO implement me
	panic("implement me")
}

func (a *articleService) Detail(ctx context.Context, id int64) (domain.Article, error) {
	// TODO implement me
	panic("implement me")
}

func (a *articleService) List(ctx context.Context, uid, page, size int64) ([]domain.Article, error) {
	return a.repo.GetByAuthorId(ctx, uid, page, size)
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, d domain.Article) (int64, error) {
	if d.Id <= 0 {
		// 新增
		return a.repo.Create(ctx, d)
	}
	// 修改
	// 有一种做法：先根据id查询文档，不存在就报错，存在就更新，相当于查询数据库两次，性能会比较差
	// 另一种做法就是直接更新，根据 update 语句返回的受影响行数来判断是否更新成功，性能会比较好
	return d.Id, a.repo.Update(ctx, d)
}
