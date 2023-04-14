package machine

import (
	"context"
	"github.com/jinzhu/copier"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMachineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMachineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMachineLogic {
	return &GetMachineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMachineLogic) GetMachine(req *types.GetMachineIdReq) (*types.GetMachineRsp, error) {
	var resp types.GetMachineRsp
	machine, _ := l.svcCtx.MachineModel.FindOne(l.ctx, req.Id)
	copier.Copy(&resp, &machine)

	return &resp, nil
}
