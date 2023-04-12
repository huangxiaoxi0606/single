package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MachineModel = (*customMachineModel)(nil)

type (
	// MachineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMachineModel.
	MachineModel interface {
		machineModel
	}

	customMachineModel struct {
		*defaultMachineModel
	}
)

// NewMachineModel returns a model for the database table.
func NewMachineModel(conn sqlx.SqlConn, c cache.CacheConf) MachineModel {
	return &customMachineModel{
		defaultMachineModel: newMachineModel(conn, c),
	}
}
