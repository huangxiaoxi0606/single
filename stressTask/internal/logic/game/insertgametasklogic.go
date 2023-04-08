package game

import (
	"context"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
	"single/stressTask/model"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type InsertGameTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertGameTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertGameTaskLogic {
	return &InsertGameTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertGameTaskLogic) InsertGameTask(req *types.InsertTaskGameReq) error {
	task := new(model.Task)
	copier.Copy(&task, &req)
	game := new(model.Game)
	copier.Copy(&game, &req)

	return nil
}
