package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"single/stressTask/internal/config"
	"single/stressTask/model"
)

type ServiceContext struct {
	Config           config.Config
	GameModel        model.GameModel
	TaskModel        model.TaskModel
	TaskMonitorModel model.TaskMonitorModel
	MachineModel     model.MachineModel
	ReportModel      model.ReportModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:           c,
		GameModel:        model.NewGameModel(sqlConn, c.CacheRedis),
		TaskModel:        model.NewTaskModel(sqlConn, c.CacheRedis),
		TaskMonitorModel: model.NewTaskMonitorModel(sqlConn, c.CacheRedis),
		MachineModel:     model.NewMachineModel(sqlConn, c.CacheRedis),
		ReportModel:      model.NewReportModel(sqlConn, c.CacheRedis),
	}
}
