package machine

import (
	"context"
	"github.com/jinzhu/copier"
	"single/common/paginate"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMachineListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMachineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMachineListLogic {
	return &GetMachineListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMachineListLogic) GetMachineList(req *types.GetMachineListReq) (*types.GetMachineListRsp, error) {
	var resp = &types.GetMachineListRsp{}
	count, _ := l.svcCtx.MachineModel.CountAllNotDelete(l.ctx, req.Name)
	newPaginate := paginate.NewPaginate(int(req.Page), int(req.PageSize), int(count))

	machine, _ := l.svcCtx.MachineModel.FindAllNotDelete(l.ctx, req.Name, int64(newPaginate.Offset), int64(newPaginate.Limit))
	copier.Copy(&resp.List, &machine)
	resp.Page = int64(newPaginate.Page)
	resp.PageSize = int64(newPaginate.Limit)
	resp.Total = int64(newPaginate.TotalResults)

	return resp, nil
}
