package client

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	"math/rand"
)

// 修改webook internal内部interaction服务的rpc调用
// 修改ioc组装代码

// InteractionServiceClient 将本地rpc调用、真实rpc调用，装饰成一个interaction client。
type InteractionServiceClient struct {
	local  intrv1.InteractionServiceClient
	remote intrv1.InteractionServiceClient

	threshold *atomicx.Value[int32]
}

func NewInteractionServiceClient(local intrv1.InteractionServiceClient, remote intrv1.InteractionServiceClient, threshold int32) *InteractionServiceClient {
	return &InteractionServiceClient{
		threshold: atomicx.NewValueOf[int32](threshold),
		local:     local,
		remote:    remote,
	}
}

func (client *InteractionServiceClient) View(ctx context.Context, in *intrv1.ViewReq, opts ...grpc.CallOption) (*intrv1.ViewResp, error) {
	return client.selectClient().View(ctx, in, opts...)
}

func (client *InteractionServiceClient) Like(ctx context.Context, in *intrv1.LikeReq, opts ...grpc.CallOption) (*intrv1.LikeResp, error) {
	return client.selectClient().Like(ctx, in, opts...)
}

func (client *InteractionServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeReq, opts ...grpc.CallOption) (*intrv1.CancelLikeResp, error) {
	return client.selectClient().CancelLike(ctx, in, opts...)
}

func (client *InteractionServiceClient) Favorite(ctx context.Context, in *intrv1.FavoriteReq, opts ...grpc.CallOption) (*intrv1.FavoriteResp, error) {
	return client.selectClient().Favorite(ctx, in, opts...)
}

func (client *InteractionServiceClient) Get(ctx context.Context, in *intrv1.GetReq, opts ...grpc.CallOption) (*intrv1.GetResp, error) {
	return client.selectClient().Get(ctx, in, opts...)
}

func (client *InteractionServiceClient) Liked(ctx context.Context, in *intrv1.LikedReq, opts ...grpc.CallOption) (*intrv1.LikedResp, error) {
	return client.selectClient().Liked(ctx, in, opts...)
}

func (client *InteractionServiceClient) Collected(ctx context.Context, in *intrv1.CollectedReq, opts ...grpc.CallOption) (*intrv1.CollectedResp, error) {
	return client.selectClient().Collected(ctx, in, opts...)
}

func (client *InteractionServiceClient) GetByIDs(ctx context.Context, in *intrv1.GetByIDsReq, opts ...grpc.CallOption) (*intrv1.GetByIDsResp, error) {
	return client.selectClient().GetByIDs(ctx, in, opts...)
}

func (client *InteractionServiceClient) selectClient() intrv1.InteractionServiceClient {
	num := rand.Int31n(100)
	if num < client.threshold.Load() {
		return client.remote
	}
	return client.local
}

func (client *InteractionServiceClient) UpdateThreshold(val int32) {
	client.threshold.Store(val)
}
