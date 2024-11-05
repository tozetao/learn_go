package web

import "testing"

type Article struct {
	Id      int64
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`
	// 作者
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}

// 测试发布
func TestArticleHandler_Publish(t *testing.T) {
	//testCases := []struct{
	//	name string
	//
	//	// 要提前准备的数据
	//	before func(t *testing.T)
	//
	//	// 验证并删除的数据
	//	after func(t *testing.T)
	//
	//	// 构造请求，直接使用req。也就是说放弃测试Bind的异常分支
	//	req Article
	//}
}
