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
	id, err := l.svcCtx.MachineModel.FindIdByIP(l.ctx, req.OuternetIp, req.InnernetIp)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.FindIdByIP is fail"), "l.svcCtx.MachineModel.FindIdByIP is fail err is : %+v", err)
	}
	if id > 0 && id != req.Id {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.FindIdByIP ip is exist"), "l.svcCtx.MachineModel.FindIdByIP ip is exist, req.OuternetIp is : %+v", req.OuternetIp)
	}

	machine := new(model.Machine)
	copier.Copy(&machine, &req)
	machine.UpdateTime = time.Now()
	err = l.svcCtx.MachineModel.Update(l.ctx, machine)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.Update is fail"), "l.svcCtx.MachineModel.Update is fail err is : %+v", err)
	}
	return nil
}
