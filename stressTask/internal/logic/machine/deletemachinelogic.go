package machine

import (
	"context"
	"github.com/pkg/errors"
	"single/common/xerr"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMachineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMachineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMachineLogic {
	return &DeleteMachineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMachineLogic) DeleteMachine(req *types.DeleteMachineReq) error {
	err := l.svcCtx.MachineModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.Delete is fail"), "l.svcCtx.MachineModel.Delete is fail err is : %+v", err)
	}
	return nil

	return nil
}
