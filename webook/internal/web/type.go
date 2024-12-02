package web

import "learn_go/webook/internal/domain"

const ArticleLike = 1
const ArticleUnlike = 0

type ArticleVO struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	CTime   string `json:"c_time"`
	UTime   string `json:"u_time"`
}

type FavoriteReq struct {
	ArticleID int64 `json:"article_id"`
	//收藏夹id
	FavoriteID int64 `json:"favorite_id"`
}

type LikeReq struct {
	ArticleID int64 `json:"article_id"`
	Action    int
}

type ListReq struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ArticleReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{ID: uid},
	}
}
