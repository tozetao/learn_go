package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"learn_go/webook/interaction/domain"
	"strconv"
	"time"
)

var (
	//go:embed lua/interaction.lua
	script string

	likeCntField     = "like_cnt"
	readCntField     = "read_cnt"
	favoriteCntField = "favorite_cnt"
)

// InteractionCache 使用hash来存储文章的交互信息
type InteractionCache interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	// IncrLikeCnt IncrLikeCntIfPresent
	IncrLikeCnt(ctx context.Context, biz string, bizID int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizID int64) error

	IncrFavoriteCnt(ctx context.Context, biz string, bizID int64) error
	DecrFavoriteCnt(ctx context.Context, biz string, bizID int64) error

	Get(ctx context.Context, biz string, bizID int64) (domain.Interaction, error)
	Set(ctx context.Context, biz string, bizID int64, interaction domain.Interaction) error
}

func (cache *interactionCache) Get(ctx context.Context, biz string, bizID int64) (domain.Interaction, error) {
	inter := domain.Interaction{}
	key := cache.key(biz, bizID)

	res, err := cache.cmd.HGetAll(ctx, key).Result()
	if err != nil {
		return inter, err
	}

	// 该对象可能被别人意外设置了，需要记录错误日志。
	if len(res) == 0 {
		return inter, ErrKeyNotExist
	}
	inter.Likes, _ = strconv.ParseInt(res[likeCntField], 10, 64)
	inter.Views, _ = strconv.ParseInt(res[readCntField], 10, 64)
	inter.Favorites, _ = strconv.ParseInt(res[favoriteCntField], 10, 64)
	return inter, nil
}

func (cache *interactionCache) Set(ctx context.Context, biz string, bizID int64, interaction domain.Interaction) error {
	key := cache.key(biz, bizID)

	err := cache.cmd.HSet(ctx, key,
		likeCntField, interaction.Likes,
		readCntField, interaction.Views,
		favoriteCntField, interaction.Favorites).Err()
	if err != nil {
		return err
	}
	return cache.cmd.Expire(ctx, key, time.Minute*10).Err()
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
	_, err := cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{readCntField, 1}).Result()
	return err
	// 我们不关注lua脚本执行的结果，脚本中已经判定了只有key存在值才会自增
	//fmt.Printf("incr read_cnt: %v\n", res)
}

func (cache *interactionCache) IncrLikeCnt(ctx context.Context, biz string, bizID int64) error {
	return cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{likeCntField, 1}).Err()
}

func (cache *interactionCache) DecrLikeCnt(ctx context.Context, biz string, bizID int64) error {
	return cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{likeCntField, -1}).Err()
}

func (cache *interactionCache) IncrFavoriteCnt(ctx context.Context, biz string, bizID int64) error {
	return cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{favoriteCntField, 1}).Err()
}

func (cache *interactionCache) DecrFavoriteCnt(ctx context.Context, biz string, bizID int64) error {
	return cache.cmd.Eval(ctx, script, []string{cache.key(biz, bizID)}, []any{favoriteCntField, -1}).Err()
}
