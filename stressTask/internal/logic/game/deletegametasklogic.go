package game

import (
	"context"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGameTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGameTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGameTaskLogic {
	return &DeleteGameTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGameTaskLogic) DeleteGameTask(req *types.DeleteTaskGameReq) error {
	err := l.svcCtx.TaskModel.Delete(l.ctx, req.Id)
	if err != nil {
		return err
	}
	// todo: add your logic here and delete this line

	return nil
}
