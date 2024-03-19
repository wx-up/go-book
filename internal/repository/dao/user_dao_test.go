package dao

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDAO_FindByEmail(t *testing.T) {
}

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user model.User

		wantErr error
		wantId  int64
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				result := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(result)
				return mockDB
			},
			wantId:  3,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: mockDB,
				// SELECT VERSION();
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 不需要 ping
				DisableAutomaticPing: true,
				// 默认情况下gorm 即时单执行一个 insert 也会先 begin 开启一个事务，执行完成之后 commit
				// mysql 其实默认是 autocommit 但是其他的一些数据库可能不是的，所以 gorm 提供了这种配置项
				SkipDefaultTransaction: true,
			})
			d := NewGORMUserDAO(db)
			lastId, err := d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantId, lastId)
		})
	}
}
