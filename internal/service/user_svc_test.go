package service

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"

	"github.com/wx-up/go-book/internal/domain"

	repomocks "github.com/wx-up/go-book/internal/repository/mocks"
	"go.uber.org/mock/gomock"

	"github.com/wx-up/go-book/internal/repository"
)

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		ctx  context.Context
		obj  domain.User

		wantObj domain.User
		wantErr error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "1453085314@qq.com").
					Return(domain.User{
						Email:    "1453085314@qq.com",
						Password: "$2a$10$0tCLtxrhPZ33ynwTNKKWlO3yQvJAYmiDL61my73i9quacnCtPf1Pq",
						Id:       1,
					}, nil)
				return repo
			},
			obj: domain.User{
				Email:    "1453085314@qq.com",
				Password: "12345",
			},

			wantObj: domain.User{
				Email:    "1453085314@qq.com",
				Password: "$2a$10$0tCLtxrhPZ33ynwTNKKWlO3yQvJAYmiDL61my73i9quacnCtPf1Pq",
				Id:       1,
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil)
			res, err := svc.Login(context.Background(), tc.obj)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantObj, res)
		})
	}
}

func Test(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	fmt.Println(err)
	fmt.Println(string(res))
}
