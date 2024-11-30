package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed "./lua/interaction.lua"
	script string
)

// InteractionCache 使用hash来存储文章的交互信息
type InteractionCache interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	IncrLikeCnt(ctx context.Context, biz string, bizID int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizID int64) error
}

type interactionCache struct {
	cmd redis.Cmdable
}

func NewInteractionCache(cmd redis.Cmdable) InteractionCache {
	return &interactionCache{
		cmd: cmd,
	}
}

func (cache *interactionCache) key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}

// IncrReadCnt lua脚本逻辑：只有在key存在的情况下戏赠。key存在意味着有人访问了该文章，缓存中已经载入数据库中该文章的阅读量，因此可以在自增加1。
func (cache *interactionCache) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	res, err := cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{"read_cnt", 1}).Result()
	// 我们不关注lua脚本执行的结果，脚本中已经判定了只有key存在值才会自增
	fmt.Printf("incr read_cnt: %v\n", res)
	return err
}

func (cache *interactionCache) IncrLikeCnt(ctx context.Context, biz string, bizID int64) error {
	//TODO implement me
	panic("implement me")
}

func (cache *interactionCache) DecrLikeCnt(ctx context.Context, biz string, bizID int64) error {
	//TODO implement me
	panic("implement me")
}
