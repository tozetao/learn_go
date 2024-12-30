package integration

type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	//// 作者
	AuthorID int64 `json:"author_id"`
	Ctime    int64 `json:"c_time"`
	Utime    int64 `json:"u_time"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
