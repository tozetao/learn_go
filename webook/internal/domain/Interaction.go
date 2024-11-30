package domain

import "time"

type Interaction struct {
	ID    int64
	Biz   string
	BizID int64 `json:"biz_id"`

	UTime time.Time `json:"u_time"`
	CTime time.Time `json:"c_time"`

	Views     int64
	Likes     int64
	Favorites int64
}

type ArticleLike struct {
	ID int64

	Uid   int64
	Biz   string
	BizID int64

	UTime time.Time `json:"u_time"`
	CTime time.Time `json:"c_time"`
}
