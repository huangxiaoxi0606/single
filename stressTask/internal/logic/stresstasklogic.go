package logic

import (
	"context"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StressTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStressTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StressTaskLogic {
	return &StressTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StressTaskLogic) StressTask(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
