package domain

import "time"

// Interaction 资源交互对象
type Interaction struct {
	ID    int64
	Biz   string
	BizID int64 `json:"biz_id"`

	UTime time.Time `json:"u_time"`
	CTime time.Time `json:"c_time"`

	Views     int64
	Likes     int64
	Favorites int64

	Liked     bool
	Collected bool
}

// UserLike 用户点赞对象
type UserLike struct {
	ID int64

	Uid   int64
	Biz   string
	BizID int64

	UTime time.Time `json:"u_time"`
	CTime time.Time `json:"c_time"`
}

func (l UserLike) Liked() bool {
	return l.Uid > 0 && l.BizID > 0
}

type UserFavorite struct {
	ID int64 `json:"id" gorm:"primaryKey,autoIncrement"`

	Uid int64 `json:"uid" gorm:"uniqueIndex:idx_uid_biz_type_id"`

	Biz   string `json:"biz" gorm:"type:varchar(128);uniqueIndex:idx_uid_biz_type_id"`
	BizID int64  `json:"biz_id" gorm:"uniqueIndex:idx_uid_biz_type_id"`

	// 收藏夹id是唯一的，自己有索引
	FavoriteID int64 `json:"favorite_id" gorm:"index"`

	UTime time.Time `json:"u_time" gorm:"column:u_time"`
	CTime time.Time `json:"c_time" gorm:"column:c_time"`
}

func (f UserFavorite) Collected() bool {
	return f.Uid > 0 && f.BizID > 0
}
