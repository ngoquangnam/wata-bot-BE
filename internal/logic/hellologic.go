package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
)

type HelloLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloLogic) Hello(req *types.HelloReq) (resp *types.HelloResp, err error) {
	name := req.Name
	if name == "" {
		name = "World"
	}
	
	return &types.HelloResp{
		Message: "Hello, " + name + "!",
	}, nil
}

