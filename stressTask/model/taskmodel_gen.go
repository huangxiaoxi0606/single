// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	taskFieldNames          = builder.RawFieldNames(&Task{})
	taskRows                = strings.Join(taskFieldNames, ",")
	taskRowsExpectAutoSet   = strings.Join(stringx.Remove(taskFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	taskRowsWithPlaceHolder = strings.Join(stringx.Remove(taskFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheHhxTaskIdPrefix = "cache:hhx:task:id:"
)

type (
	taskModel interface {
		Insert(ctx context.Context, data *Task) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Task, error)
		Update(ctx context.Context, data *Task) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTaskModel struct {
		sqlc.CachedConn
		table string
	}

	Task struct {
		Id            int64     `db:"id"`
		Name          string    `db:"name"`           // 压测任务名字
		Desc          string    `db:"desc"`           // 压测描述
		TotalNum      int64     `db:"total_num"`      // 总人数
		Pace          int64     `db:"pace"`           // 步调
		RunTime       int64     `db:"run_time"`       // 执行时间（s）
		OpNickname    string    `db:"op_nickname"`    // 执行人
		MachineConfig string    `db:"machine_config"` // 机器设置
		MaxRps        int64     `db:"max_rps"`        // 最大rps【阈值】
		IsDeleted     int64     `db:"is_deleted"`     // 0不删除，1已删除
		CreateTime    time.Time `db:"create_time"`    // 创建时间
		UpdateTime    time.Time `db:"update_time"`    // 更新时间
		DeleteTime    time.Time `db:"delete_time"`    // 删除时间
		LastDoTime    time.Time `db:"last_do_time"`   // 最后一次执行时间
		TaskFlag      int64     `db:"task_flag"`      // 类型1game2pressure
	}
)

func newTaskModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTaskModel {
	return &defaultTaskModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`task`",
	}
}

func (m *defaultTaskModel) Delete(ctx context.Context, id int64) error {
	hhxTaskIdKey := fmt.Sprintf("%s%v", cacheHhxTaskIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, hhxTaskIdKey)
	return err
}

func (m *defaultTaskModel) FindOne(ctx context.Context, id int64) (*Task, error) {
	hhxTaskIdKey := fmt.Sprintf("%s%v", cacheHhxTaskIdPrefix, id)
	var resp Task
	err := m.QueryRowCtx(ctx, &resp, hhxTaskIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", taskRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTaskModel) Insert(ctx context.Context, data *Task) (sql.Result, error) {
	hhxTaskIdKey := fmt.Sprintf("%s%v", cacheHhxTaskIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, taskRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.Desc, data.TotalNum, data.Pace, data.RunTime, data.OpNickname, data.MachineConfig, data.MaxRps, data.IsDeleted, data.DeleteTime, data.LastDoTime, data.TaskFlag)
	}, hhxTaskIdKey)
	return ret, err
}

func (m *defaultTaskModel) Update(ctx context.Context, data *Task) error {
	hhxTaskIdKey := fmt.Sprintf("%s%v", cacheHhxTaskIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, taskRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Name, data.Desc, data.TotalNum, data.Pace, data.RunTime, data.OpNickname, data.MachineConfig, data.MaxRps, data.IsDeleted, data.DeleteTime, data.LastDoTime, data.TaskFlag, data.Id)
	}, hhxTaskIdKey)
	return err
}

func (m *defaultTaskModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheHhxTaskIdPrefix, primary)
}

func (m *defaultTaskModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", taskRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTaskModel) tableName() string {
	return m.table
}
