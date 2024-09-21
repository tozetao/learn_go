package domain

import (
	"time"
)

/*
Domain：领域对象，根据业务边界所抽象出来的对象。
Service: 领域服务，代表一个业务的完成流程。
Repository：代表领域对象的存储。
*/

type User struct {
	ID       int64
	Email    string
	Password string
	Phone    string

	Nickname string
	AboutMe  string
	Birthday time.Time
	Ctime    time.Time

	WechatInfo WechatInfo
}
