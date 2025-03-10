package integration

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"Neo/Workplace/goland/src/GeekGo/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name:   "发送成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:13800138000").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody: `{
				"phone":"13800138000"
			}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				//这个手机号已经有一个验证码了
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				_, err := rdb.Set(ctx, "phone_code:login:13800138000", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:13800138000").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是 6 位,没有被覆盖
				assert.Equal(t, "123456", val)
			},
			reqBody: `{
				"phone":"13800138000"
			}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送验证码太频繁, 1分钟后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				//这个手机号已经有一个验证码了，但是没有过期时间
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				_, err := rdb.Set(ctx, "phone_code:login:13800138000", "123456",
					0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:13800138000").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是 6 位,没有被覆盖
				assert.Equal(t, "123456", val)
			},
			reqBody: `{
				"phone":"13800138000"
			}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码为空",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `{
				"phone":""
			}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
		{
			name: "数据格式错误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `{
				"phone":"",
			}`,
			wantCode: http.StatusBadRequest,
			//wantBody: web.Result{
			//	Code: 4,
			//	Msg:  "输入有误",
			//},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest("POST",
				"/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			//t.Logf("%+v", resp)

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)

			if tc.wantCode != http.StatusOK {
				return
			}

			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}
