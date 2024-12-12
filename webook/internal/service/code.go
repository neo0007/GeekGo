package service

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (svc CodeService) Send(ctx context.Context,
//区别业务场景
	biz string,
	phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.smsSvc.Send(ctx, codeTplId, []string{code}, code)
}

func (svc CodeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)

}

func (svc *CodeService) generateCode() string {
	// 6位数，num 在 0-999999之间，包括 0 和 999999
	num := rand.Intn(1000000)
	//不够 6 位数，加上前导 0
	return fmt.Sprintf("%06d", num)

}
