// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

var (
	reportFieldNames          = builder.RawFieldNames(&Report{})
	reportRows                = strings.Join(reportFieldNames, ",")
	reportRowsExpectAutoSet   = strings.Join(stringx.Remove(reportFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	reportRowsWithPlaceHolder = strings.Join(stringx.Remove(reportFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheHhxReportIdPrefix = "cache:hhx:report:id:"
)

type (
	reportModel interface {
		Insert(ctx context.Context, data *Report) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Report, error)
		Update(ctx context.Context, data *Report) error
		Delete(ctx context.Context, id int64) error
		UpdatePlatReportMonitorPng(ctx context.Context, id int64, ReportMonitorPng sql.NullString) error
		UpdateReportFlow(ctx context.Context, id int64, ReceiveFlow, TransmitFlow string) error
	}

	defaultReportModel struct {
		sqlc.CachedConn
		table string
	}

	Report struct {
		Id               int64          `db:"id"`                 // id
		TaskId           int64          `db:"task_id"`            // 任务id
		Result           int64          `db:"result"`             // 执行结果
		Executor         string         `db:"executor"`           // 操作人
		TotalNum         int64          `db:"total_num"`          // 总人数
		Pace             int64          `db:"pace"`               // 步调
		ReportDesc       string         `db:"report_desc"`        // 报告描述
		MachineConfig    string         `db:"machine_config"`     // 发压机配置
		DistData         sql.NullString `db:"dist_data"`          // 请求的统计总数据
		ReqData          sql.NullString `db:"req_data"`           // 请求的统计时间数据
		ErrData          sql.NullString `db:"err_data"`           // 错误数据
		RunTime          int64          `db:"run_time"`           // 执行时间(s)
		StartTime        sql.NullTime   `db:"start_time"`         // 开始执行时间
		EndTime          sql.NullTime   `db:"end_time"`           // 执行结束时间
		IsDelete         int64          `db:"is_delete"`          // 是否删除 0不删除1删除
		UpdateTime       sql.NullTime   `db:"update_time"`        // 更新时间
		ReceiveFlow      string         `db:"receive_flow"`       // 收到流量
		TransmitFlow     string         `db:"transmit_flow"`      // 发送流量
		ExecutorType     int64          `db:"executor_type"`      // 任务执行类型1业务执行2定时任务执行
		MaxTps           int64          `db:"max_tps"`            // 最大tps
		ReportMonitorPng sql.NullString `db:"report_monitor_png"` // 报告监控信息截图
		TaskType         int64          `db:"task_type"`          // 类型1game2pressure
		Name             string         `db:"name"`               // 报告名称
	}
)

func newReportModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultReportModel {
	return &defaultReportModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`report`",
	}
}

func (m *defaultReportModel) Delete(ctx context.Context, id int64) error {
	hhxReportIdKey := fmt.Sprintf("%s%v", cacheHhxReportIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, hhxReportIdKey)
	return err
}

func (m *defaultReportModel) FindOne(ctx context.Context, id int64) (*Report, error) {
	hhxReportIdKey := fmt.Sprintf("%s%v", cacheHhxReportIdPrefix, id)
	var resp Report
	err := m.QueryRowCtx(ctx, &resp, hhxReportIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", reportRows, m.table)
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

func (m *defaultReportModel) Insert(ctx context.Context, data *Report) (sql.Result, error) {
	hhxReportIdKey := fmt.Sprintf("%s%v", cacheHhxReportIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, reportRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.TaskId, data.Result, data.Executor, data.TotalNum, data.Pace, data.ReportDesc, data.MachineConfig, data.DistData, data.ReqData, data.ErrData, data.RunTime, data.StartTime, data.EndTime, data.IsDelete, data.ReceiveFlow, data.TransmitFlow, data.ExecutorType, data.MaxTps, data.ReportMonitorPng, data.TaskType, data.Name)
	}, hhxReportIdKey)
	return ret, err
}

func (m *defaultReportModel) Update(ctx context.Context, data *Report) error {
	hhxReportIdKey := fmt.Sprintf("%s%v", cacheHhxReportIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, reportRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.TaskId, data.Result, data.Executor, data.TotalNum, data.Pace, data.ReportDesc, data.MachineConfig, data.DistData, data.ReqData, data.ErrData, data.RunTime, data.StartTime, data.EndTime, data.IsDelete, data.ReceiveFlow, data.TransmitFlow, data.ExecutorType, data.MaxTps, data.ReportMonitorPng, data.TaskType, data.Name, data.Id)
	}, hhxReportIdKey)
	return err
}

func (m *defaultReportModel) UpdatePlatReportMonitorPng(ctx context.Context, id int64, ReportMonitorPng sql.NullString) error {
	hhxReportIdKey := fmt.Sprintf("%sUpdatePlatReportMonitorPng%v", cacheHhxReportIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set`ReportMonitorPng` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, ReportMonitorPng, id)
	}, hhxReportIdKey)
	return err
}

func (m *defaultReportModel) UpdateReportFlow(ctx context.Context, id int64, ReceiveFlow, TransmitFlow string) error {
	hhxReportIdKey := fmt.Sprintf("%UpdateReportFlow%v", cacheHhxReportIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set`ReceiveFlow` = ? ,`TransmitFlow` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, ReceiveFlow, TransmitFlow, id)
	}, hhxReportIdKey)
	return err
}

func (m *defaultReportModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheHhxReportIdPrefix, primary)
}

func (m *defaultReportModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", reportRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultReportModel) tableName() string {
	return m.table
}
