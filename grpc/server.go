package grpc

import "context"

type Server struct {
	// UnsafeUserServiceServer接口
	// 每次grpc生成的UnimplementedUserServiceServer实现了UnsafeUserServiceServer接口，同时也实现了UserService接口。
	// 嵌入该结构体，当idl service新增了接口，即使你代码没有去实现这些新增的接口，你的程序也不会编译报错。
	UnimplementedUserServiceServer

	Name string
}

func (s *Server) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		User: &User{
			Id:   123,
			Name: "from" + s.Name,
		},
	}, nil
}
