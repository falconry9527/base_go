// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.5
// Source: demo.proto

package democlient

import (
	"context"

	"base_go/04_zero/02_grpc/demo/demo"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Request  = demo.Request
	Response = demo.Response

	Demo interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	}

	defaultDemo struct {
		cli zrpc.Client
	}
)

func NewDemo(cli zrpc.Client) Demo {
	return &defaultDemo{
		cli: cli,
	}
}

func (m *defaultDemo) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := demo.NewDemoClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}
