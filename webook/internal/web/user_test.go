package web

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	svcmocks "Neo/Workplace/goland/src/GeekGo/webook/internal/service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "signup success",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					Return(nil)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"12345678",
			"confirmPassword":"12345678"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "signup successfully",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"12345678",
		}
		`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq",
			"password":"12345678",
			"confirmPassword":"12345678"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对！",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"12345678",
			"confirmPassword":"14345678"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致！",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"123478",
			"confirmPassword":"123478"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于 8 位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"12345678",
			"confirmPassword":"12345678"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "邮箱或手机号码冲突了！",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					Return(errors.New("system error"))
				return userSvc
			},
			reqBody: `
		{
			"email":"123@qq.com",
			"password":"12345678",
			"confirmPassword":"12345678"
		}
		`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()

			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest("POST",
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			//t.Logf("%+v", resp)

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, resp.Body.String())
			}
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := svcmocks.NewMockUserService(ctrl)

	userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	err := userSvc.Signup(context.Background(), domain.User{Email: "123@qq"})
	//t.Log(err)
	require.Equal(t, "mock error", err.Error())
}
