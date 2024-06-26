package service

import (
	"context"
	domain2 "github.com/wx-up/go-book/interactive/domain"
	"github.com/wx-up/go-book/interactive/service"
	"testing"
	"time"

	"github.com/wx-up/go-book/internal/domain"

	svcmocks "github.com/wx-up/go-book/internal/service/mocks"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestBatchRankingService_TopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService)
		wantErr error
		wanRes  []int64
	}{
		{
			name: "计算成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				// hacknews 模型受时间的影响，为了方便测试，这里将时间固定
				// 这样子的话，点赞数越大，score 越大，文章排名越靠前
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), int64(0), int64(3)).Return(
					[]domain.Article{
						{
							Id:         1,
							CreateTime: now,
							UpdateTime: now,
						},
						{
							Id:         2,
							CreateTime: now,
							UpdateTime: now,
						},
						{
							Id:         3,
							CreateTime: now,
							UpdateTime: now,
						},
					}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), int64(3), int64(3)).Return(nil, nil)
				intSvc := svcmocks.NewMockInteractiveService(ctrl)
				intSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2, 3}).Return(
					[]domain2.Interactive{
						{
							BizId:   1,
							LikeCnt: 10,
						},
						{
							BizId:   2,
							LikeCnt: 100,
						},
						{
							BizId:   3,
							LikeCnt: 30,
						},
					},
					nil)
				intSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).Return(nil, nil)
				return artSvc, nil
			},
			wantErr: nil,
			wanRes:  []int64{2, 3, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, intSvc := tc.mock(ctrl)
			svc := NewArticleRankingService(artSvc, intSvc, nil)
			svc.batchSize = 3
			svc.n = 3
			svc.scoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			ids, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wanRes, ids)
		})
	}
}

func Ha(artSvc ArticleService) {
	for i := 0; i < 3; i++ {
		artSvc.Detail(context.Background(), 0)
	}
}

func TestH(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	artSvc := svcmocks.NewMockArticleService(ctrl)
	artSvc.EXPECT().Detail(context.Background(), int64(0)).Return(domain.Article{}, nil)
	Ha(artSvc)
}
