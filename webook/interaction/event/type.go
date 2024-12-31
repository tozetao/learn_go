package event

const TopicArticleReadEvent = "article_read"

type ReadEvent struct {
	Uid       int64 `json:"uid"`
	ArticleID int64 `json:"article_id"`
}
