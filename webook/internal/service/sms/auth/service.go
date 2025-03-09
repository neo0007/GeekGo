package auth

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

// Send 发送，其中 biz 必须是线下申请的一个代表业务方的 token
func (s *SMSService) Send(ctx context.Context, biz string,
	args []string, numbers ...string) error {

	var tc Claims
	// 如果这里能解析成功，说明就是对应的业务方
	// 这里直接加密再解密，实际使用中把加密算法放在须授权独立模块
	token, err := jwt.ParseWithClaims(s.bizTokenString("300100"),
		&tc, func(t *jwt.Token) (interface{}, error) {
			return s.key, nil
		})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}

	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

func NewAuthSMSService(svc sms.Service, key []byte) sms.Service {
	return &SMSService{
		svc: svc,
		key: key,
	}
}

func (s *SMSService) bizTokenString(tpl string) string {
	claims := Claims{
		Tpl: tpl,
	}
	bizToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bizTokenString, err := bizToken.SignedString(s.key)
	if err != nil {
		zap.L().Warn("bizToken err", zap.Error(err))
	}
	return bizTokenString
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
