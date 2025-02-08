package failover

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	smsmocks "Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sync/atomic"
	"testing"
)

func Test_timeoutFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.Service
		threshold int32
		cnt       int32
		wantCnt   int32
		wantIdx   int32
		wantErr   error
	}{
		{
			name: "触发了切换，切换之后成功了",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)

				return []sms.Service{svc0, svc1}
			},
			cnt:       3,
			threshold: 3,
			// 重置了
			wantCnt: 0,
			wantIdx: 1,
			wantErr: nil,
		},
		{
			name: "触发了切换，切换之后依旧超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)

				return []sms.Service{svc0, svc1}
			},
			cnt:       3,
			threshold: 3,
			// 重置了
			wantCnt: 1,
			wantIdx: 1,
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "触发了切换，切换之后失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))

				return []sms.Service{svc0, svc1}
			},
			cnt:       3,
			threshold: 3,
			// 重置了
			wantCnt: 0,
			wantIdx: 1,
			wantErr: errors.New("发送失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), 0, tc.cnt, tc.threshold)
			err := svc.Send(context.Background(), "mytpl", []string{"123"}, "13901390000")
			assert.Equal(t, tc.wantErr, err)
			timeoutSvc, ok := svc.(*timeoutFailoverSMSService)
			assert.True(t, ok, "svc should be of type *timeoutFailoverSMSService")
			assert.Equal(t, tc.wantCnt, atomic.LoadInt32(&timeoutSvc.cnt), "cnt mismatch")
			assert.Equal(t, tc.wantIdx, atomic.LoadInt32(&timeoutSvc.idx), "idx mismatch")
		})
	}
}
