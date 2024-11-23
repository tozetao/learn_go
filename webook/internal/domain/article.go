package domain

import "time"

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus

	CTime time.Time
	UTime time.Time
}

type Author struct {
	ID   int64
	Name string
}

type ArticleStatus int8

func (a ArticleStatus) ToInt8() int8 {
	return int8(a)
}

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)
