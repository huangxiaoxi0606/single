package game

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"single/common/xerr"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
	"single/stressTask/model"

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

// InsertGameTask 创建任务
func (l *InsertGameTaskLogic) InsertGameTask(req *types.InsertTaskGameReq) error {
	task := new(model.Task)
	copier.Copy(&task, &req)
	game := new(model.Game)
	copier.Copy(&game, &req)
	MachineConfigByte, err := json.Marshal(&req.MachineConfig)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("MachineConfig data is not supported"), "MachineConfig is wrong json req: %+v", req)
	}
	task.MachineConfig = string(MachineConfigByte)
	//插入task
	taskData, err := l.svcCtx.TaskModel.Insert(l.ctx, task)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskModel.Insert is fail"), "l.svcCtx.TaskModel.Insert is fail err is : %+v", err)
	}
	//获取刚保存的taskId
	taskId, _ := taskData.LastInsertId()
	game.TaskId = taskId
	//插入game
	_, err = l.svcCtx.GameModel.Insert(l.ctx, game)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.GameModel.Insert is fail"), "l.svcCtx.GameModel.Insert is fail err is : %+v", err)
	}
	//解析monitor并插入
	if len(req.MonitorConfig) > 0 { //可优化为批量插入
		for _, val := range req.MonitorConfig {
			monitor := new(model.TaskMonitor)
			copier.Copy(&monitor, &val)
			monitor.TaskId = taskId
			_, err = l.svcCtx.TaskMonitorModel.Insert(l.ctx, monitor)
			if err != nil {
				return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.TaskMonitorModel.Insert is fail"), "l.svcCtx.TaskMonitorModel.Insert is fail err is : %+v", err)
			}
		}
	}
	return nil
}
