package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server

	Addr string
}

func (server *Server) Start() error {
	lis, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	return server.Serve(lis)
}
