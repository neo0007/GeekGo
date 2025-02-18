package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/oauth2/wechat"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
)

func InitWechatService() wechat.Service {
	appId := "This_is_for_Wechat_appId"
	appSecret := "This_is_for_Wechat_appSecret"
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID environment variable is not set")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("WECHAT_APP_SECRET environment variable is not set")
	//}
	return wechat.NewService(appId, appSecret)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
