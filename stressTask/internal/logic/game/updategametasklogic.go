package game

import (
	"context"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGameTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameTaskLogic {
	return &UpdateGameTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGameTaskLogic) UpdateGameTask(req *types.UpdateTaskGameReq) error {
	// todo: add your logic here and delete this line

	return nil
}
