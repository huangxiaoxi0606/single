package machine

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"single/common/xerr"
	"single/stressTask/model"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InsertMachineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertMachineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertMachineLogic {
	return &InsertMachineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertMachineLogic) InsertMachine(req *types.InsertMachineReq) error {
	machine := new(model.Machine)
	copier.Copy(&machine, &req)
	b, err := l.svcCtx.MachineModel.FindExistByIP(l.ctx, machine.OuternetIp, machine.InnernetIp)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskModel.FindExistByIP is fail"), "l.svcCtx.TaskModel.FindExistByIP is fail err is : %+v", err)
	}
	if b {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskModel.FindExistByIP Ip is exist"), "l.svcCtx.TaskModel.FindExistByIP Ip is exist ip is : %+v", machine.InnernetIp)
	}
	_, err = l.svcCtx.MachineModel.Insert(l.ctx, machine)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.MachineModel.Insert is fail"), "l.svcCtx.MachineModel.Insert is fail err is : %+v", err)
	}
	return nil
}
