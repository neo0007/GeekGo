package ratelimit

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"Neo/Workplace/goland/src/GeekGo/webook/pkg/ratelimit"
	"context"
	"fmt"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 你在这里添加一些代码
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		//系统错误，即 redis 错误
		//限流：保守策略，你的下游很弱的时候
		//不限流：容错策略，你的下游很强，业务可用性要求很高
		//包一下这个错误
		return fmt.Errorf("短信服务判断是否限流出现问题: %w", err)
	}
	if limited {
		return errLimited
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 你在这里也可以添加一些新特性
	return err
}
