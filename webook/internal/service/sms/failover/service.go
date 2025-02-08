package failover

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context,
	tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		//输出日志，做好监控
		log.Println(err)
	}
	//所有服务商失败，意味着你的网络崩了
	return errors.New("发送失败，所有服务商都失败")
}

func (f *FailoverSMSService) SendV1(ctx context.Context,
	tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 调用者这设置的超时时间到了
			// 调用者主动取消
			return err
		default:
			// 其他情况会走到这里，要打印日志监控
			log.Println(err)
		}
	}
	return errors.New("发送失败，所有服务商都失败")
}
