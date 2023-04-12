package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TaskMonitorModel = (*customTaskMonitorModel)(nil)

type (
	// TaskMonitorModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskMonitorModel.
	TaskMonitorModel interface {
		taskMonitorModel
	}

	customTaskMonitorModel struct {
		*defaultTaskMonitorModel
	}
)

// NewTaskMonitorModel returns a model for the database table.
func NewTaskMonitorModel(conn sqlx.SqlConn, c cache.CacheConf) TaskMonitorModel {
	return &customTaskMonitorModel{
		defaultTaskMonitorModel: newTaskMonitorModel(conn, c),
	}
}
