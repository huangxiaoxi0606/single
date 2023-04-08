package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PressureModel = (*customPressureModel)(nil)

type (
	// PressureModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPressureModel.
	PressureModel interface {
		pressureModel
	}

	customPressureModel struct {
		*defaultPressureModel
	}
)

// NewPressureModel returns a model for the database table.
func NewPressureModel(conn sqlx.SqlConn, c cache.CacheConf) PressureModel {
	return &customPressureModel{
		defaultPressureModel: newPressureModel(conn, c),
	}
}
