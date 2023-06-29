/*
@Time : 2023/4/19 20:54
@Author : Hhx06
@File : util
@Description:
@Software: GoLand
*/

package report

import (
	"encoding/json"
	"github.com/pkg/errors"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
	"single/stressTask/model"
	"strconv"
)

//验证发压机是否可用,并返回主控机和发压机
func (l *RunTaskLogic) verifyMachine(svcCtx *svc.ServiceContext, task *model.Task) (master *model.Machine, slaves []*model.Machine, err error) {
	var (
		MachineConfig types.StressMachine
	)

	err = json.Unmarshal([]byte(task.MachineConfig), &MachineConfig)
	if err != nil {
		return
	}
	if MachineConfig.Number < 2 || len(MachineConfig.MachineIds) < 1 {
		err = errors.New("Machine count is < 2")
		return
	}
	if MachineConfig.MasterMachineId == 0 {
		err = errors.New("MachineConfig.MasterMachineId is empty")
		return
	}
	//判断主控机是否被占用
	master, err = svcCtx.MachineModel.FindOne(l.ctx, MachineConfig.MasterMachineId)
	if err != nil {
		return
	}
	if master.IsDelete == 1 || master.UseFlag == 1 || master.WorkingFlag == 1 {
		err = errors.New("MachineConfig.MasterMachineId is  unavailable,MachineConfig.MasterMachineId:" + strconv.FormatInt(MachineConfig.MasterMachineId, 10))
		return
	}
	//判断发压机是否被占用
	for _, v := range MachineConfig.MachineIds {
		var slave *model.Machine
		slave, err = svcCtx.MachineModel.FindOne(l.ctx, v)
		if err != nil {
			return
		}
		if slave.IsDelete == 1 || slave.UseFlag == 1 || slave.WorkingFlag == 1 {
			err = errors.New("MachineConfig.MasterMachineIds is  unavailable,MachineConfig.slave:" + strconv.FormatInt(v, 10))
			return
		}
		slaves = append(slaves, slave)
	}
	return

}
