package game

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"single/common/xerr"
	"single/stressTask/model"
	"time"

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
	task := new(model.Task)
	copier.Copy(&task, &req)
	game := new(model.Game)
	copier.Copy(&game, &req)
	MachineConfigByte, err := json.Marshal(&req.MachineConfig)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("MachineConfig data is not supported"), "MachineConfig is wrong json req: %+v", req)
	}
	task.UpdateTime = time.Now()
	task.DeleteTime = time.Now()
	task.LastDoTime = time.Now()
	task.MachineConfig = string(MachineConfigByte)
	err = l.svcCtx.TaskModel.Update(l.ctx, task)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskModel.Update is fail"), "l.svcCtx.TaskModel.Update is fail err is : %+v", err)
	}
	game.TaskId = task.Id
	err = l.svcCtx.GameModel.UpdateByTaskId(l.ctx, game)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.GameModel.Update is fail"), "l.svcCtx.GameModel.Update is fail err is : %+v", err)
	}
	//删除已有的
	err = l.svcCtx.TaskMonitorModel.DeleteByTaskId(l.ctx, task.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.GameModel.DeleteByTaskId is fail"), "l.svcCtx.GameModel.DeleteByTaskId is fail err is : %+v", err)
	}
	//解析monitor并插入
	if len(req.MonitorConfig) > 0 { //可优化为批量插入
		for _, val := range req.MonitorConfig {
			monitor := new(model.TaskMonitor)
			copier.Copy(&monitor, &val)
			monitor.TaskId = task.Id
			_, err = l.svcCtx.TaskMonitorModel.Insert(l.ctx, monitor)
			if err != nil {
				return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskMonitorModel.Insert is fail"), "l.svcCtx.TaskMonitorModel.Insert is fail err is : %+v", err)
			}
		}
	}
	return nil
}
