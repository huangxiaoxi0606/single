package machine

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"single/common/xerr"
	"single/stressTask/model"
	"time"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMachineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMachineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMachineLogic {
	return &UpdateMachineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMachineLogic) UpdateMachine(req *types.UpdateMachineReq) error {
	machine := new(model.Machine)
	copier.Copy(&machine, &req)
	machine.UpdateTime = time.Now()
	err := l.svcCtx.MachineModel.Update(l.ctx, machine)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.Update is fail"), "l.svcCtx.MachineModel.Update is fail err is : %+v", err)
	}
	return nil
}
