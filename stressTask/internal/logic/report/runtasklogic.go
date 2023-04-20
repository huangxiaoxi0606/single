package report

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

type RunTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunTaskLogic {
	return &RunTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunTaskLogic) RunTask(req *types.InsertReportReq) error {
	var task *model.Task
	task, _ = l.svcCtx.TaskModel.FindOne(l.ctx, req.TaskId)
	if task.IsDeleted == 1 { //任务已经被删除
		return errors.Wrapf(xerr.NewErrMsg("task is deleted"), "task is deleted req.TaskId: %+v", req.TaskId)
	}
	if req.TotalNum < 1 {
		return errors.Wrapf(xerr.NewErrMsg("req.TotalNum is wrong"), "req.TotalNum is wrong req.TotalNum: %+v", req.TotalNum)
	}

	err := l.verifyMachine(l.svcCtx, task) //验证发压机是否可用
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("verifyMachine is wrong"), "verifyMachine is wrong err: %+v", err)
	}
	//先创建对应的测试报告
	report := new(model.Report)
	copier.Copy(&report, &req)
	report.MachineConfig = task.MachineConfig
	report.TaskType = task.TaskFlag
	report.Result = 0
	_, err = l.svcCtx.ReportModel.Insert(l.ctx, report)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("create report is failed"), "create report is failed err: %+v", err)
	}

	//var report

	return nil
}
