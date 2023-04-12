package game

import (
	"context"
	"github.com/jinzhu/copier"
	"single/common/paginate"
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

func (l *GetGameTaskListLogic) GetGameTaskList(req *types.GetGameTaskListReq) (*types.GetGameTaskListRsp, error) {
	var resp = &types.GetGameTaskListRsp{}
	count, _ := l.svcCtx.TaskModel.CountAllNotDelete(l.ctx, req.Name)
	newPaginate := paginate.NewPaginate(int(req.Page), int(req.PageSize), int(count))

	task, _ := l.svcCtx.TaskModel.FindAllNotDelete(l.ctx, req.Name, int64(newPaginate.Offset), int64(newPaginate.Limit))
	//var data = new([]types.GetGameTaskRsp)
	copier.Copy(&resp.List, &task)
	resp.Page = int64(newPaginate.Page)
	resp.PageSize = int64(newPaginate.Limit)
	resp.Total = int64(newPaginate.TotalResults)

	return resp, nil
}
