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

//验证发压机是否可用
func (l *RunTaskLogic) verifyMachine(svcCtx *svc.ServiceContext, task *model.Task) error {
	var MachineConfig types.StressMachine
	err := json.Unmarshal([]byte(task.MachineConfig), &MachineConfig)
	if err != nil {
		return errors.New("json.Unmarshal failed")
	}
	if MachineConfig.Number < 2 || len(MachineConfig.MachineIds) < 1 {
		return errors.New("Machine count is < 2")
	}
	if MachineConfig.MasterMachineId == 0 {
		return errors.New("MachineConfig.MasterMachineId is empty")
	}
	//判断主控机是否被占用
	master, err := svcCtx.MachineModel.FindOne(l.ctx, MachineConfig.MasterMachineId)
	if err != nil {
		return err
	}
	if master.IsDelete == 1 || master.UseFlag == 1 || master.WorkingFlag == 1 {
		return errors.New("MachineConfig.MasterMachineId is  unavailable,MachineConfig.MasterMachineId:" + strconv.FormatInt(MachineConfig.MasterMachineId, 10))
	}
	//判断发压机是否被占用
	for _, v := range MachineConfig.MachineIds {
		slave, err := svcCtx.MachineModel.FindOne(l.ctx, v)
		if err != nil {
			return err
		}
		if slave.IsDelete == 1 || slave.UseFlag == 1 || slave.WorkingFlag == 1 {
			return errors.New("MachineConfig.MasterMachineIds is  unavailable,MachineConfig.slave:" + strconv.FormatInt(v, 10))
		}
	}
	return nil

}
