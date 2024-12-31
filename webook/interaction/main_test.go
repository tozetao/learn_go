package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cc, err := grpc.Dial("localhost:8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := intrv1.NewInteractionServiceClient(cc)

	resp, err := client.View(ctx, &intrv1.ViewReq{
		BizId: 1,
		Biz:   "article",
	})

	assert.NoError(t, err)
	t.Logf("resp: %v", resp)
	t.Logf("resp1: %v", &intrv1.ViewResp{})
	//assert.Equal(t, &intrv1.ViewResp{}, resp)
}
