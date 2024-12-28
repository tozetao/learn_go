package grpc

//func TestServer(t *testing.T) {
//	gs := grpc.NewServer()
//	us := &Server{
//		Name: "lisi-server",
//	}
//	RegisterUserServiceServer(gs, us)
//
//	l, err := net.Listen("tcp", ":8090")
//	require.NoError(t, err)
//
//	err = gs.Serve(l)
//	t.Log(err)
//}
//
//func TestClient(t *testing.T) {
//	cc, err := grpc.Dial("localhost:8090",
//		grpc.WithTransportCredentials(insecure.NewCredentials()))
//	require.NoError(t, err)
//
//	client := NewUserServiceClient(cc)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	resp, err := client.GetById(ctx, &GetByIdReq{
//		Id: 123,
//	})
//	assert.NoError(t, err)
//	t.Log(resp.User)
//}
