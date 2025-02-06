package memory

import (
	mySMS "Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"fmt"
)

type Service struct {
}

func NewService() mySMS.Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
