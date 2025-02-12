package retryable

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"errors"
)

type Service struct {
	svc sms.Service
	// 重试
	retryMax int
}

func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, biz, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, biz, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}
func NewService(svc sms.Service, retryMax int) sms.Service {
	return &Service{
		svc:      svc,
		retryMax: retryMax,
	}
}
