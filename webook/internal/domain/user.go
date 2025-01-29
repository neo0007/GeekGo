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
	Ctime    time.Time
	Utime    time.Time
}
