package grpc

import (
	"context"
	"learn_go/webook/api/proto/gen/intr"
	"learn_go/webook/interaction/domain"
	"learn_go/webook/interaction/service"
)

type InteractionServiceServer struct {
	svc service.InteractionService

	intrv1.UnimplementedInteractionServiceServer
}

func NewInteractionServiceServer(svc service.InteractionService) *InteractionServiceServer {
	return &InteractionServiceServer{
		svc: svc,
	}
}

func (server *InteractionServiceServer) View(ctx context.Context, req *intrv1.ViewReq) (*intrv1.ViewResp, error) {
	err := server.svc.View(ctx, req.GetBiz(), req.GetBizId())
	return &intrv1.ViewResp{}, err
}

func (server *InteractionServiceServer) Like(ctx context.Context, req *intrv1.LikeReq) (*intrv1.LikeResp, error) {
	err := server.svc.Like(ctx, req.GetUid(), req.GetBiz(), req.GetBizId())
	return &intrv1.LikeResp{}, err
}

func (server *InteractionServiceServer) CancelLike(ctx context.Context, req *intrv1.CancelLikeReq) (*intrv1.CancelLikeResp, error) {
	err := server.svc.CancelLike(ctx, req.GetUid(), req.GetBiz(), req.GetBizId())
	return &intrv1.CancelLikeResp{}, err
}

func (server *InteractionServiceServer) Favorite(ctx context.Context, req *intrv1.FavoriteReq) (*intrv1.FavoriteResp, error) {
	err := server.svc.Favorite(ctx, req.GetUid(), req.GetFavoriteId(), req.GetBiz(), req.GetBizId())
	return &intrv1.FavoriteResp{}, err
}

func (server *InteractionServiceServer) Get(ctx context.Context, req *intrv1.GetReq) (*intrv1.GetResp, error) {
	inter, err := server.svc.Get(ctx, req.GetUid(), req.GetBiz(), req.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResp{
		Inter: server.toDTO(inter),
	}, nil
}

func (server *InteractionServiceServer) Liked(ctx context.Context, req *intrv1.LikedReq) (*intrv1.LikedResp, error) {
	liked, err := server.svc.Liked(ctx, req.GetUid(), req.GetBiz(), req.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.LikedResp{Liked: liked}, nil
}

func (server *InteractionServiceServer) Collected(ctx context.Context, req *intrv1.CollectedReq) (*intrv1.CollectedResp, error) {
	collected, err := server.svc.Collected(ctx, req.GetUid(), req.GetBiz(), req.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectedResp{Collected: collected}, nil
}

func (server *InteractionServiceServer) GetByIDs(ctx context.Context, req *intrv1.GetByIDsReq) (*intrv1.GetByIDsResp, error) {
	m, err := server.svc.GetByIDs(ctx, req.GetBiz(), req.GetBizIds())
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*intrv1.Interaction, len(m))
	for i, inter := range m {
		res[i] = server.toDTO(inter)
	}
	return &intrv1.GetByIDsResp{
		Inters: res,
	}, nil
}

// data transfer object
func (server *InteractionServiceServer) toDTO(inter domain.Interaction) *intrv1.Interaction {
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
