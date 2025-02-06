package web

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	svcmocks "Neo/Workplace/goland/src/GeekGo/webook/internal/service/mocks"
	"bytes"
	"context"
	"encoding/json"
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

func TestUserHandler_LoginSMS(t *testing.T) {
	const biz = "login"
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		ctx      context.Context
		reqBody  string
		wantJSON Result
	}{
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Id:    int64(1),
						Phone: "13912345678",
					}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 5,
				Msg:  "验证码校验通过",
			},
		},
		{
			name: "bind error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().FindOrCreate(context.Background(), gomock.Any()).
				//	Return(domain.User{
				//		Id:    int64(1),
				//		Phone: "13912345678",
				//	}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				//codeSvc.EXPECT().Verify(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
				//	Return(true, nil)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456",
		}
		`,
			wantJSON: Result{
				Code: 5,
				Msg:  "bind error",
			},
		},
		{
			name: "验证次数太多",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
				//	Return(domain.User{
				//		Id:    int64(1),
				//		Phone: "13912345678",
				//	}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, service.ErrCodeVerifyTooManyTimes)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 4,
				Msg:  "验证次数太多，请重新发送验证码",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
				//	Return(domain.User{
				//		Id:    int64(1),
				//		Phone: "13912345678",
				//	}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, errors.New("system error"))
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证码有误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
				//	Return(domain.User{
				//		Id:    int64(1),
				//		Phone: "13912345678",
				//	}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, nil)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 4,
				Msg:  "验证码有误",
			},
		},
		{
			name: "数据库系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("DB system error"))

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 5,
				Msg:  "DB系统错误",
			},
		},
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Id:    int64(1),
						Phone: "13912345678",
					}, nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)
				return userSvc, codeSvc
			},
			ctx: context.Background(),
			reqBody: `
		{
			"phone":"13912345678",
			"code":"123456"
		}
		`,
			wantJSON: Result{
				Code: 5,
				Msg:  "验证码校验通过",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			h := NewUserHandler(tc.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest("POST",
				"/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			//t.Logf("%+v", resp)

			server.ServeHTTP(resp, req)
			var respJSON Result
			err = json.Unmarshal(resp.Body.Bytes(), &respJSON)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantJSON, respJSON)
		})
	}
}
