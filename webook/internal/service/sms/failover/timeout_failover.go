package failover

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"errors"
	"sync/atomic"
)

type timeoutFailoverSMSService struct {
	// 你的服务商
	svcs []sms.Service
	idx  int32
	// 连续超时的个数
	cnt int32

	//阈值，连续超过这个数字就要切换
	threshold int32
}

func (t *timeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		// 这里要切换新的下标，往后挪一个
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 成功的往后挪了一位
			idx++
			atomic.StoreInt32(&t.cnt, 0)
		} else {
			//	这里出现并发，别人换成功了
			idx = atomic.LoadInt32(&t.idx)
		}
		svc := t.svcs[idx]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			// 没有任何错误，重置计数器
			atomic.StoreInt32(&t.cnt, 0)
		case context.DeadlineExceeded:
			atomic.AddInt32(&t.cnt, 1)
		default:
			// 如果是别的异常的话，保持不变
			err = errors.New("发送失败")
		}
		return err
	}
	return nil
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, idx int32, cnt int32, threshold int32) sms.Service {
	return &timeoutFailoverSMSService{
		svcs:      svcs,
		idx:       idx,
		cnt:       cnt,
		threshold: threshold,
	}
}
