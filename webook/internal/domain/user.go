package domain

import "time"

// User 领域对象，是 DDD 中的 entity
// BO Business Object
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	Birthday string
	AboutMe  string
	// 不要组合，万一你将来还有 DingDingInfo, 里面可能有同名字段
	WechatInfo WechatInfo
	Ctime      time.Time
	Utime      time.Time
}
