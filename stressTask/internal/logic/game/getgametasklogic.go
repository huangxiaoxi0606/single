package game

import (
	"context"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGameTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameTaskLogic {
	return &GetGameTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGameTaskLogic) GetGameTask(req *types.GetGameTaskIdReq) (resp *types.GetGameTaskRsp, err error) {
	// todo: add your logic here and delete this line

	return
}
