package integration

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"github.com/wx-up/go-book/interactive/integration/startup"
	"github.com/wx-up/go-book/interactive/repository/dao/model"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"testing"
	"time"
)

type InteractiveTestSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
}

func (s *InteractiveTestSuite) SetupSuite() {
	s.db = startup.InitTestMysql()
	s.rdb = startup.InitTestRedis()
}

func (s *InteractiveTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `interactives`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_like_bizs`").Error
	assert.NoError(s.T(), err)
}

func (s *InteractiveTestSuite) TestIncrReadCnt() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64

		wantErr  error
		wantResp *inter.IncrReadCntResponse
	}{
		{
			// DB 和缓存都有数据
			name: "增加成功,db和redis",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(model.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
					UpdateTime: 7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 7)
				data.UpdateTime = 0
				assert.Equal(t, model.Interactive{
					Id:    1,
					Biz:   "test",
					BizId: 2,
					// +1 之后
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
				}, data)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			wantResp: &inter.IncrReadCntResponse{},
		},
		{
			// DB 有数据，缓存没有数据
			name: "增加成功,db有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(model.Interactive{
					Id:         2,
					Biz:        "test",
					BizId:      3,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
					UpdateTime: 7,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("id = ?", 2).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 7)
				data.UpdateTime = 0
				assert.Equal(t, model.Interactive{
					Id:    2,
					Biz:   "test",
					BizId: 3,
					// +1 之后
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    3,
			wantResp: &inter.IncrReadCntResponse{},
		},
		{
			name:   "增加成功-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("biz_id = ? AND biz = ?", 4, "test").First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 0)
				assert.True(t, data.CreateTime > 0)
				assert.True(t, data.Id > 0)
				data.UpdateTime = 0
				data.CreateTime = 0
				data.Id = 0
				assert.Equal(t, model.Interactive{
					Biz:     "test",
					BizId:   4,
					ReadCnt: 1,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    4,
			wantResp: &inter.IncrReadCntResponse{},
		},
	}

	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			res, err := svc.IncrReadCnt(context.Background(), &inter.IncrReadCntRequest{
				Biz:   tc.biz,
				BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				tc.after(t)
				return
			}
			assert.Equal(t, tc.wantResp, res)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr error
	}{
		{
			name: "点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(model.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
					UpdateTime: 7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 7)
				data.UpdateTime = 0
				assert.Equal(t, model.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    6,
					CreateTime: 6,
				}, data)

				var likeBiz model.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 2, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.CreateTime > 0)
				assert.True(t, likeBiz.UpdateTime > 0)
				likeBiz.Id = 0
				likeBiz.CreateTime = 0
				likeBiz.UpdateTime = 0
				assert.Equal(t, model.UserLikeBiz{
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Status: 1,
				}, likeBiz)

				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:   "test",
			bizId: 2,
			uid:   123,
		},
		{
			name:   "点赞-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?",
					"test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 0)
				assert.True(t, data.LikeCnt > 0)
				assert.True(t, data.Id > 0)
				data.UpdateTime = 0
				data.CreateTime = 0
				data.Id = 0
				assert.Equal(t, model.Interactive{
					Biz:     "test",
					BizId:   3,
					LikeCnt: 1,
				}, data)

				var likeBiz model.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 3, 124).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.CreateTime > 0)
				assert.True(t, likeBiz.UpdateTime > 0)
				likeBiz.Id = 0
				likeBiz.CreateTime = 0
				likeBiz.UpdateTime = 0
				assert.Equal(t, model.UserLikeBiz{
					Biz:    "test",
					BizId:  3,
					Uid:    124,
					Status: 1,
				}, likeBiz)

				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:   "test",
			bizId: 3,
			uid:   124,
		},
	}

	svc := startup.InitInteractiveService()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := svc.Like(context.Background(), tc.biz, tc.bizId, tc.uid)
			assert.NoError(t, err)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestDislike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr error
	}{
		{
			name: "取消点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(model.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CreateTime: 6,
					UpdateTime: 7,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(model.UserLikeBiz{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					Uid:        123,
					CreateTime: 6,
					UpdateTime: 7,
					Status:     1,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data model.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UpdateTime > 7)
				data.UpdateTime = 0
				assert.Equal(t, model.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    4,
					CreateTime: 6,
				}, data)

				var likeBiz model.UserLikeBiz
				err = s.db.Where("id = ?", 1).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.UpdateTime > 7)
				likeBiz.UpdateTime = 0
				assert.Equal(t, model.UserLikeBiz{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					Uid:        123,
					CreateTime: 6,
					Status:     0,
				}, likeBiz)

				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:   "test",
			bizId: 2,
			uid:   123,
		},
	}

	svc := startup.InitInteractiveService()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := svc.CancelLike(context.Background(), tc.biz, tc.bizId, tc.uid)
			assert.NoError(t, err)
			tc.after(t)
		})
	}
}

func TestInteractiveService(t *testing.T) {
	suite.Run(t, &InteractiveTestSuite{})
}
