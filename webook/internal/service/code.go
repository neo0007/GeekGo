package service

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = repository.ErrUnknownForCode
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type DesignCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &DesignCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}

}

func (svc *DesignCodeService) Send(ctx context.Context,
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

func (svc *DesignCodeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)

}

func (svc *DesignCodeService) generateCode() string {
	// 6位数，num 在 0-999999之间，包括 0 和 999999
	num := rand.Intn(1000000)
	//不够 6 位数，加上前导 0
	return fmt.Sprintf("%06d", num)

}
