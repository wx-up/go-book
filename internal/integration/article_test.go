package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wx-up/go-book/internal/web/jwt"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"gorm.io/gorm"

	"github.com/wx-up/go-book/internal/integration/startup"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

// ArticleTestSuite 另一种测试用例的组织形式
type ArticleTestSuite struct {
	suite.Suite
	engine *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	s.engine = gin.Default()
	// 因为目前测试的对象是 ArticleHandler 不需要将整个项目启动起来，即调用 InitWebService
	// 并且整个项目启动流程也比较复杂，而且还有很多中间件鉴权等操作，为了方便测试
	// 这里只注册 ArticleHandler 的路由，以及相关的依赖
	s.engine.Use(func(context *gin.Context) {
		context.Set("claims", jwt.UserClaim{
			Uid: 10,
		})
	})
	ah := startup.CreateArticleHandler()
	ah.RegisterRoutes(s.engine)
	s.db = startup.InitTestMysql()
}

func (s *ArticleTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE articles")
}

type Article struct {
	ID      int
	Title   string
	Content string
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (s *ArticleTestSuite) TestSave() {
	testCases := []struct {
		name string
		req  any

		// 准备数据
		before func(t *testing.T)
		// 清除数据
		after func(t *testing.T)

		wantCode int
		wantData Result[map[string]int]
	}{
		{
			name:     "【参数有误】非法json",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			req:      123,
			wantCode: http.StatusBadRequest,
			wantData: Result[map[string]int]{
				Code: -1,
				Msg:  "参数错误",
			},
		},
		{
			name:   "【参数有误】标题为空",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			req: Article{
				Title:   "",
				Content: "test content",
			},
			wantCode: http.StatusBadRequest,
			wantData: Result[map[string]int]{
				Code: -1,
				Msg:  "参数错误",
			},
		},
		{
			name:   "新增文章成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				var obj model.Article
				err := s.db.Where("title = ?", "test title").First(&obj).Error
				require.NoError(t, err)
				assert.True(t, obj.CreateTime > 0)
				assert.True(t, obj.UpdateTime > 0)
				obj.CreateTime = 0
				obj.UpdateTime = 0
				assert.Equal(t, obj, model.Article{
					Id:       1,
					Title:    "test title",
					Content:  "test content",
					AuthorId: 10,
				})
			},
			req: Article{
				Title:   "test title",
				Content: "test content",
			},
			wantCode: http.StatusOK,
			wantData: Result[map[string]int]{
				Code: 0,
				Msg:  "保存成功",
				Data: map[string]int{
					"id": 1,
				},
			},
		},
		{
			name: "修改已经存在的文章",
			before: func(t *testing.T) {
				s.db.Create(&model.Article{
					Id:         2,
					Title:      "这是标题",
					Content:    "这是内容",
					AuthorId:   10,
					CreateTime: 1000,
					UpdateTime: 1000,
				})
			},
			after: func(t *testing.T) {
				var obj model.Article
				err := s.db.Where("id = ?", 2).First(&obj).Error
				require.NoError(t, err)
				// 更新操作：创建时间不变
				assert.Equal(t, obj.CreateTime, int64(1000))
				// 更新时间更新为当前时间，直接断言相等有点困难，这里只判断更新时间大于之前的更新时间
				assert.True(t, obj.UpdateTime > 1000)
				obj.UpdateTime = 0
				assert.Equal(t, obj, model.Article{
					Id:         2,
					Title:      "这是修改了的标题",
					Content:    "这是修改了的内容",
					AuthorId:   10,
					CreateTime: 1000,
				})
			},
			req: Article{
				ID:      2,
				Title:   "这是修改了的标题",
				Content: "这是修改了的内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[map[string]int]{
				Code: 0,
				Msg:  "保存成功",
				Data: map[string]int{
					"id": 2,
				},
			},
		},
		{
			name: "【安全策略】修改不存在的文章",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			req: Article{
				ID:      3,
				Title:   "这是修改了的标题",
				Content: "这是修改了的内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[map[string]int]{
				Code: 1,
				Msg:  "参数错误",
			},
		},
		{
			name: "【安全策略】修改别人的文章",
			before: func(t *testing.T) {
				s.db.Create(&model.Article{
					Id:         4,
					Title:      "这是标题",
					Content:    "这是内容",
					AuthorId:   4,
					CreateTime: 1000,
					UpdateTime: 1000,
				})
			},
			after: func(t *testing.T) {
				var obj model.Article
				err := s.db.Where("id = ?", 4).First(&obj).Error
				require.NoError(t, err)
				assert.Equal(t, obj.CreateTime, int64(1000))
				assert.Equal(t, obj.UpdateTime, int64(1000))
				assert.Equal(t, obj, model.Article{
					Id:         4,
					Title:      "这是标题",
					Content:    "这是内容",
					AuthorId:   4,
					CreateTime: 1000,
					UpdateTime: 1000,
				})
			},
			req: Article{
				ID:      4,
				Title:   "这是修改了的标题",
				Content: "这是修改了的内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[map[string]int]{
				Code: 1,
				Msg:  "参数错误",
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			data, err := json.Marshal(tc.req)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/save", bytes.NewReader(data))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.engine.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}
			var res Result[map[string]int]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, res)
		})
	}
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
