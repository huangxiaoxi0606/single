package game

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGameTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameTaskLogic {
	return &GetGameTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetGameTask 获取任务
func (l *GetGameTaskLogic) GetGameTask(req *types.GetGameTaskIdReq) (*types.GetGameTaskRsp, error) {
	var resp types.GetGameTaskRsp
	monitors, _ := l.svcCtx.TaskMonitorModel.FindAllByTaskId(l.ctx, req.Id)
	var ms []types.GameMonitorConfig
	if len(monitors) > 0 {
		for _, value := range monitors {
			ms = append(ms, types.GameMonitorConfig{
				Name:    value.Name,
				Url:     value.Url,
				Pwd:     value.Pwd,
				Type:    value.Type,
				Account: value.Account,
			})
		}
	}
	resp.MonitorConfig = ms
	task, _ := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	copier.Copy(&resp, &task)
	resp.CreateTime = task.CreateTime.String()
	resp.LastDoTime = task.LastDoTime.String()
	json.Unmarshal([]byte(task.MachineConfig), resp.MachineConfig)
	game, _ := l.svcCtx.GameModel.FindOneByTaskId(l.ctx, req.Id)
	copier.Copy(&resp, &game)
	resp.Id = req.Id
	return &resp, nil
}
