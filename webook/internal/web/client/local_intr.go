package client

import (
	"context"
	"google.golang.org/grpc"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	"learn_go/webook/interaction/domain"
	"learn_go/webook/interaction/service"
)

// InteractionServiceAdapter 将本地的interaction服务伪装成rpc client。
type InteractionServiceAdapter struct {
	svc service.InteractionService
}

func NewInteractionServiceAdapter(svc service.InteractionService) *InteractionServiceAdapter {
	return &InteractionServiceAdapter{svc: svc}
}

func (adapter *InteractionServiceAdapter) View(ctx context.Context, in *intrv1.ViewReq, opts ...grpc.CallOption) (*intrv1.ViewResp, error) {
	err := adapter.svc.View(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.ViewResp{}, err
}

func (adapter *InteractionServiceAdapter) Like(ctx context.Context, in *intrv1.LikeReq, opts ...grpc.CallOption) (*intrv1.LikeResp, error) {
	err := adapter.svc.Like(ctx, in.GetUid(), in.GetBiz(), in.GetBizId())
	return &intrv1.LikeResp{}, err
}

func (adapter *InteractionServiceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeReq, opts ...grpc.CallOption) (*intrv1.CancelLikeResp, error) {
	err := adapter.svc.CancelLike(ctx, in.GetUid(), in.GetBiz(), in.GetBizId())
	return &intrv1.CancelLikeResp{}, err
}

func (adapter *InteractionServiceAdapter) Favorite(ctx context.Context, in *intrv1.FavoriteReq, opts ...grpc.CallOption) (*intrv1.FavoriteResp, error) {
	err := adapter.svc.Favorite(ctx, in.GetUid(), in.GetFavoriteId(), in.GetBiz(), in.GetBizId())
	return &intrv1.FavoriteResp{}, err
}

func (adapter *InteractionServiceAdapter) Get(ctx context.Context, in *intrv1.GetReq, opts ...grpc.CallOption) (*intrv1.GetResp, error) {
	inter, err := adapter.svc.Get(ctx, in.GetUid(), in.GetBiz(), in.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResp{
		Inter: adapter.toDTO(inter),
	}, nil
}

func (adapter *InteractionServiceAdapter) Liked(ctx context.Context, in *intrv1.LikedReq, opts ...grpc.CallOption) (*intrv1.LikedResp, error) {
	liked, err := adapter.svc.Liked(ctx, in.GetUid(), in.GetBiz(), in.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.LikedResp{Liked: liked}, nil
}

func (adapter *InteractionServiceAdapter) Collected(ctx context.Context, in *intrv1.CollectedReq, opts ...grpc.CallOption) (*intrv1.CollectedResp, error) {
	collected, err := adapter.svc.Collected(ctx, in.GetUid(), in.GetBiz(), in.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectedResp{Collected: collected}, nil
}

func (adapter *InteractionServiceAdapter) GetByIDs(ctx context.Context, in *intrv1.GetByIDsReq, opts ...grpc.CallOption) (*intrv1.GetByIDsResp, error) {
	m, err := adapter.svc.GetByIDs(ctx, in.GetBiz(), in.GetBizIds())
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*intrv1.Interaction, len(m))
	for i, inter := range m {
		res[i] = adapter.toDTO(inter)
	}
	return &intrv1.GetByIDsResp{
		Inters: res,
	}, nil
}

// data transfer object
func (adapter *InteractionServiceAdapter) toDTO(inter domain.Interaction) *intrv1.Interaction {
	return &intrv1.Interaction{
		Id:        inter.ID,
		Biz:       inter.Biz,
		BizId:     inter.BizID,
		CTime:     inter.CTime.UnixMilli(),
		UTime:     inter.UTime.UnixMilli(),
		Views:     inter.Views,
		Likes:     inter.Likes,
		Favorites: inter.Favorites,

		Liked:     inter.Liked,
		Collected: inter.Collected,
	}
}
