package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
