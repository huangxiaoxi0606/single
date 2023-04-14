package game

import (
	"context"
	"github.com/pkg/errors"
	"single/common/xerr"

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
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskModel.Delete is fail"), "l.svcCtx.TaskModel.Delete is fail err is : %+v", err)
	}
	return nil
}
