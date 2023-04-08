package game

import (
	"context"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameTaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGameTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameTaskListLogic {
	return &GetGameTaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGameTaskListLogic) GetGameTaskList(req *types.GetGameTaskListReq) (resp *types.GetGameTaskListRsp, err error) {
	// todo: add your logic here and delete this line

	return
}
