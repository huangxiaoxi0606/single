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
	machineFieldNames          = builder.RawFieldNames(&Machine{})
	machineRows                = strings.Join(machineFieldNames, ",")
	machineRowsExpectAutoSet   = strings.Join(stringx.Remove(machineFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	machineRowsWithPlaceHolder = strings.Join(stringx.Remove(machineFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheHhxMachineIdPrefix = "cache:hhx:machine:id:"
)

type (
	machineModel interface {
		Insert(ctx context.Context, data *Machine) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Machine, error)
		Update(ctx context.Context, data *Machine) error
		Delete(ctx context.Context, id int64) error
		FindAllNotDelete(ctx context.Context, name string, offset, limit int64) ([]*Machine, error)
		CountAllNotDelete(ctx context.Context, name string) (int64, error)
		FindExistByIP(ctx context.Context, outerNetIp, innerNetIp string) (bool, error)
		FindIdByIP(ctx context.Context, outerNetIp, innerNetIp string) (int64, error)
	}

	defaultMachineModel struct {
		sqlc.CachedConn
		table string
	}

	Machine struct {
		Id           int64         `db:"id"`            // id
		Name         string        `db:"name"`          // 别名
		OuternetIp   string        `db:"outernet_ip"`   // 外网ip
		InnernetIp   string        `db:"innernet_ip"`   // 内网ip
		CpuCores     int64         `db:"cpu_cores"`     // cpu核数
		Port         int64         `db:"port"`          // 端口
		RamSize      sql.NullInt64 `db:"ram_size"`      // 内存
		RootAccount  string        `db:"root_account"`  // root账号
		RootPassword string        `db:"root_password"` // root密码
		UseFlag      int64         `db:"use_flag"`      // 是否已启用0启用1禁用
		NickName     string        `db:"nick_name"`     // 添加人名字
		WorkingFlag  int64         `db:"working_flag"`  // 是否正在工作0默认1正在工作
		CreateTime   time.Time     `db:"create_time"`   // 创建时间
		UpdateTime   time.Time     `db:"update_time"`   // 更新时间
		IsDelete     int64         `db:"is_delete"`     // 是否已删除0未删除1已删除
	}
)

func newMachineModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMachineModel {
	return &defaultMachineModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`machine`",
	}
}

func (m *defaultMachineModel) Delete(ctx context.Context, id int64) error {
	hhxMachineIdKey := fmt.Sprintf("%s%v", cacheHhxMachineIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `is_delete`= ?  where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, 1, id)
	}, hhxMachineIdKey)
	return err
}

func (m *defaultMachineModel) FindOne(ctx context.Context, id int64) (*Machine, error) {
	hhxMachineIdKey := fmt.Sprintf("%s%v", cacheHhxMachineIdPrefix, id)
	var resp Machine
	err := m.QueryRowCtx(ctx, &resp, hhxMachineIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", machineRows, m.table)
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

func (m *defaultMachineModel) Insert(ctx context.Context, data *Machine) (sql.Result, error) {
	hhxMachineIdKey := fmt.Sprintf("%s%v", cacheHhxMachineIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, machineRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.OuternetIp, data.InnernetIp, data.CpuCores, data.Port, data.RamSize, data.RootAccount, data.RootPassword, data.UseFlag, data.NickName, data.WorkingFlag, data.IsDelete)
	}, hhxMachineIdKey)
	return ret, err
}

func (m *defaultMachineModel) Update(ctx context.Context, data *Machine) error {
	hhxMachineIdKey := fmt.Sprintf("%s%v", cacheHhxMachineIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, machineRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Name, data.OuternetIp, data.InnernetIp, data.CpuCores, data.Port, data.RamSize, data.RootAccount, data.RootPassword, data.UseFlag, data.NickName, data.WorkingFlag, data.IsDelete, data.Id)
	}, hhxMachineIdKey)
	return err
}

func (m *defaultMachineModel) FindAllNotDelete(ctx context.Context, name string, offset, limit int64) ([]*Machine, error) {
	var resp []*Machine
	var err error
	if len(name) > 0 {
		err = m.QueryRowsNoCache(&resp, fmt.Sprintf("select %s from %s where `name` like ? and `is_delete`= ? order by `id` desc limit ?, ?", machineRows, m.table), name, 0, offset, limit)
	} else {
		err = m.QueryRowsNoCache(&resp, fmt.Sprintf("select %s from %s where `is_delete`= ? order by `id` desc limit ?, ?", machineRows, m.table), 0, offset, limit)
	}
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMachineModel) CountAllNotDelete(ctx context.Context, name string) (int64, error) {
	var count int64
	var err error
	if len(name) > 0 {
		err = m.QueryRowNoCache(&count, fmt.Sprintf("select count(1) from %s where `name` like ? and `is_deleted`= ?", m.table), name, 0)
	} else {
		err = m.QueryRowNoCache(&count, fmt.Sprintf("select count(1) from %s where `is_deleted`= ?", m.table), 0)
	}
	switch err {
	case nil:
		return count, nil
	default:
		return 0, err
	}
}

func (m *defaultMachineModel) FindExistByIP(ctx context.Context, outerNetIp, innerNetIp string) (bool, error) {
	var count int64
	var err error
	err = m.QueryRowNoCache(&count, fmt.Sprintf("select count(1) from %s where `outernet_ip` = ? and `innernet_ip`= ? and `is_delete` = ?", m.table), outerNetIp, innerNetIp, 0)

	switch err {
	case nil:
		if count > 0 {
			return true, nil
		}
		return false, nil
	default:
		return false, err
	}
}

func (m *defaultMachineModel) FindIdByIP(ctx context.Context, outerNetIp, innerNetIp string) (int64, error) {
	var resp Machine
	var err error
	err = m.QueryRowNoCache(&resp, fmt.Sprintf("select %s from %s where `outernet_ip` = ? and `innernet_ip`= ? and `is_delete` = ? limit 1", machineRows, m.table), outerNetIp, innerNetIp, 0)

	switch err {
	case nil:
		return resp.Id, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

func (m *defaultMachineModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheHhxMachineIdPrefix, primary)
}

func (m *defaultMachineModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", machineRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMachineModel) tableName() string {
	return m.table
}
