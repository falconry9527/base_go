package logic

import (
	"base_go/04_zero/02_grpc/demo/demo"
	"base_go/04_zero/02_grpc/demo/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *demo.Request) (*demo.Response, error) {
	// todo: add your logic here and delete this line
	return &demo.Response{
		Pong: "pong",
	}, nil
}
