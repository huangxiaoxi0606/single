package report

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"log"
	"net/http"
	"single/common"
	"single/common/aboutTps"
	"single/common/git"
	"single/common/parseData"
	"single/common/sftp"
	"single/common/ssh"
	"single/common/xerr"
	"single/stressTask/model"
	"strconv"
	"strings"
	"sync"
	"time"

	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunTaskLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	masterServerChan chan struct{}
	uploadFileChan   chan struct{}
}

type RunConfig struct {
	Pace     int64
	TotalNum int64
	RunTime  int64
	MaxTps   int64
}

func NewRunTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunTaskLogic {
	return &RunTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunTaskLogic) RunTask(req *types.InsertReportReq) error {
	var (
		task     *model.Task
		dirPath  string
		game     *model.Game
		err      error
		reportId int64
		master   *model.Machine
		slaves   []*model.Machine
	)
	task, _ = l.svcCtx.TaskModel.FindOne(l.ctx, req.TaskId)
	if task.IsDeleted == 1 { //任务已经被删除
		return errors.Wrapf(xerr.NewErrMsg("task is deleted"), "task is deleted req.TaskId: %+v", req.TaskId)
	}
	if req.TotalNum < 1 {
		return errors.Wrapf(xerr.NewErrMsg("req.TotalNum is wrong"), "req.TotalNum is wrong req.TotalNum: %+v", req.TotalNum)
	}

	master, slaves, err = l.verifyMachine(l.svcCtx, task) //验证发压机是否可用
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("verifyMachine is wrong"), "verifyMachine is wrong err: %+v", err)
	}

	err = mr.Finish(func() (err error) { //game 编译 pressure 生成
		if task.TaskFlag == 1 { //game
			game, err = l.svcCtx.GameModel.FindOneByTaskId(l.ctx, task.Id)
			if err != nil {
				return errors.Wrapf(xerr.NewErrMsg("game task is failed"), "game task is failed err: %+v", err)
			}
			dirPath, err = git.CloneDir(game.Url, l.svcCtx.Config.Git.GitlabToken, game.Branch, game.EntryDir, l.svcCtx.Config.Git.GitProjectPath, l.svcCtx.Config.Git.GitCloneUrlHead)
			if err != nil {
				logx.Errorf(" git clone err, err is %v , gitlab token is %v,  task is %v ", err, l.svcCtx.Config.Git.GitlabToken,
					task)
				return
			}
			return
		}
		return
	}, func() (err error) { //创建对应的测试报告
		report := new(model.Report)
		copier.Copy(&report, &req)
		report.MachineConfig = task.MachineConfig
		report.TaskType = task.TaskFlag
		report.Result = 0
		reportData, err := l.svcCtx.ReportModel.Insert(l.ctx, report)

		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("create report is failed"), "create report is failed err: %+v", err)
		}
		reportId, _ = reportData.LastInsertId()
		return
	})
	if err != nil {
		log.Printf("run error: %v", err)
		return errors.Wrapf(xerr.NewErrMsg("create report is failed"), "create report is failed err: %+v", err)
	}
	fmt.Println("report:", reportId)
	runConfig := &RunConfig{
		Pace:     req.Pace,
		TotalNum: req.TotalNum,
		RunTime:  req.RunTime,
		MaxTps:   req.MaxTps,
	}
	//修改机器状态

	//先创建对应的测试报告

	//var report

	return nil
}

func (l *RunTaskLogic) locustRun(task *model.Task, masterServer *model.Machine, slaveServer []*model.Machine, config *RunConfig, reportId int64, startTime, endTime time.Time) {
	var (
		runAfter      = time.After(time.Second * 300) // 如果5min后还没有上传成功，主动退出(主要兼顾海外发压机)
		err           error
		slaveRunErrCh = make(chan struct{})
		WorkerPool    = goroutine.Default()
	)

	select {
	case <-l.uploadFileChan:
		// 清除主控机进程
		err = ssh.KillOldTest(masterServer, slaveServer)
		if err != nil {
			logx.Errorf("[locustRun] kill old process err , err: %v, master: %v, slave: %v", err, masterServer, slaveServer)
			l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportId, Status: common.GameClearProcess, MachineConfig: task.MachineConfig}
			return
		}

		// 先启动主控机
		masterShell := func() {
			//server *model.Machine, serverCount int64
			err = l.MasterShell(task, config, reportId, startTime, endTime, masterServer, int64(len(slaveServer)))
			if err != nil {
				logx.Errorf("[locustRun] master run shell, err is %v , master server is %v ", err, masterServer)
				l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportId, Status: common.GameRunMaster, MachineConfig: task.MachineConfig}
			}

			WorkerPool.Submit(func() { //SaveFlow
				defer func() {
					recover()
				}()

				err = parseData.SaveFlow(l.ctx, l.svcCtx.ReportModel, l.svcCtx.Config.Flow.Url, reportId, slaveServer)
				if err != nil {
					logx.Errorf("[parseData.SaveFlow] err , err: %v, reportId: %v", err, reportId)
				}
			})

			// 执行结束后主动清理进程
			defer func(masterServer *model.Machine, slaveServer []*model.Machine) {
				err := ssh.KillOldTest(masterServer, slaveServer)
				if err != nil {
					logx.Errorf("[locustRun] kill old process err , err: %v, master: %v, slave: %v", err, masterServer, slaveServer)
					l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportId, Status: common.GameClearProcess, MachineConfig: task.MachineConfig}
				}
			}(masterServer, slaveServer)
		}
		SafeRun(masterShell) // TODO 做法待定
		// 启动slave
		// 如果 10 s 内master启动成功，任务执行失败，主动中断
		masterAfter := time.After(10 * time.Second)
		select {
		case <-l.masterServerChan:
			var slaveWg sync.WaitGroup
			for _, server := range slaveServer {
				slaveWg.Add(1)
				go func(server *model.Machine) {
					err = SlaveShell(server, masterServer.OuternetIp, strconv.FormatInt(config.MaxTps, 10))
					if err != nil {
						l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportId, Status: common.GameFail, MachineConfig: task.MachineConfig}
						logx.Errorf("[locustRun] slave run shell err, err is %v, server is %v ", err, server)
						close(slaveRunErrCh)
					}
					slaveWg.Done()
				}(server)
			}
			slaveWg.Wait()

			// todo 修改 状态为开始执行
			err = l.svcCtx.Dao.UpdateGameReportStartTime(reportId, time.Now().Unix())
			if err != nil {
				logx.Errorf("[locustRun] update report begin execute time err, reportId is %v, err is %v", reportId, err)
				//  更新报告开始执行的时间错误，无需返回, 只影响压测报告精度
			}

		case <-masterAfter:
			logx.Errorf("[locustRun] master start timeout, masterServer: %v", masterServer)
			l.svcCtx.GameStatusCh <- &dto.GameStatusParam{ReportId: reportId, Status: common.GameUploadMaster, ServerIds: serverIds}
		}
	case <-runAfter:
		logx.Errorf("[locustRun] upload file timeout, masterServer: %v", masterServer) //
		l.svcCtx.GameStatusCh <- &dto.GameStatusParam{ReportId: reportId, Status: common.GameUploadMaster, ServerIds: serverIds}
	case <-slaveRunErrCh:
		// 执行错误
		l.svcCtx.GameStatusCh <- &dto.GameStatusParam{ReportId: reportId, Status: common.GameRunSlave, ServerIds: serverIds}
	}
}

/*
   --no-web：指定无 web UI模式
   -c：起多少 locust 用户(等同于起多少 tcp 连接)
   -r：多少时间内，把上述 -c 设置的虚拟用户全部启动
   -t：脚本运行多少时间，单位
*/

func (l *RunTaskLogic) MasterShell(task *model.Task, config *RunConfig, reportID int64, startTime, endTime time.Time, server *model.Machine, serverCount int64) error {
	s := &ssh.Service{Host: server.OuternetIp, User: server.RootAccount, Pwd: server.RootPassword, Type: "password", KeyPath: "", Port: server.Port}
	shell := fmt.Sprintf("cd %s; ulimit -n 65535 ; locust --master  --no-web  --no-reset-stats --only-summary "+
		"--expect-slaves %d --logfile %d_locust.log  --master-bind-port 6667 "+
		"--csv %d -c %d -r %d -t %dm -f %s  2>&1 &",
		l.svcCtx.Config.Ssh.MasterRemotely, serverCount, reportID, reportID, config.TotalNum, config.Pace, config.RunTime, "master.py")
	sshClient, err := s.Connect()
	if err != nil || sshClient == nil {
		logx.Errorf("[MasterShell] master connect  err: %v, sshClient : %v ", err, sshClient)
		return err
	}
	// 删除普罗米修斯推送网关
	defer func() {
		l.deletePushGateway(reportID)
		sshClient.Close()
		s.Client.Close()
	}()
	if err = sshClient.Start(shell); err != nil {
		l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportID, Status: common.GameRunMaster, MachineConfig: task.MachineConfig}
		logx.Errorf("[MasterShell] master run shell is  err: %v, shell :%v", err, shell)
		return err
	}
	// 命令执行成功后通知slave开始启动
	close(l.masterServerChan)
	// 等待locust 运行结束
	if err = sshClient.Wait(); err != nil {
		l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportID, Status: common.GameFail, MachineConfig: task.MachineConfig}
		logx.Errorf("[MasterShell] master Wait shell is  err: %v", err)
		return err
	}
	// 添加报告数据
	if err = l.GetReportGame(l.ctx, reportID, startTime, endTime, common.GameSuc, server); err != nil {
		l.svcCtx.GameStatusCh <- &svc.GameStatusParam{ReportId: reportID, Status: common.GameGetReport, MachineConfig: task.MachineConfig}
		logx.Errorf("[MasterShell] get report err, err is %v", err)
		return err
	}

	defer sshClient.Close()
	return err
}

func SlaveShell(server *model.UtestStressMachine, masterIp string, maxRps string) error {
	s := &ssh.SSHService{Host: server.OuterNetIp, User: common.LINUX_ROOT, Pwd: server.RootPassword, Type: "password", KeyPath: "", Port: server.Port}
	tName := time.Now().Format("2006-01-02-15-04-05")
	shell := fmt.Sprintf("chmod  777  /data/load/slave/main;ulimit -n 100000;nohup " + "/data/load/slave/main --master-host " + masterIp + maxRps + fmt.Sprintf(" --master-port 6667 >/data/load/slave/game_%s_locust.log 2>&1 &", tName))
	sshClient, err := s.Connect()
	if err != nil || sshClient == nil {
		logx.Errorf("[SlaveShell] slave connect err, err is %v, sshClient is %v ", err, sshClient)
		return err
	}

	if err = sshClient.Run(shell); err != nil {
		logx.Errorf("[SlaveShell] slave run shell  err: %v, shell %v ", err, shell)
		return err
	}
	defer sshClient.Close()
	return err
}

//删除普罗米修斯推送网关
func (l *RunTaskLogic) deletePushGateway(reportId int64) {
	deleteUrl := fmt.Sprintf(l.svcCtx.Config.PushGateway.UrlPre+"%d", reportId)
	deleteJob, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		logx.Errorf("[deletePushGateway] delete pushGateWay err, err is %v, url is %v ", err, deleteUrl)
	}
	_, err = http.DefaultClient.Do(deleteJob)
	if err != nil {
		logx.Errorf("[deletePushGateway] deletePushGateway err: %v, url is %v ", err, deleteUrl)
	}
}

// GetReportGame 获取报告数据
func (l *RunTaskLogic) GetReportGame(ctx context.Context, reportID int64, startTime, endTime time.Time, result int64, master *model.Machine) error {
	var (
		id           = strconv.FormatInt(reportID, 10)
		c            = new(sftp.FTPClient)
		srcDistData  []*parseData.DistInfo
		distData     *parseData.DistInfo
		totalTps     float64
		maxTps       float64
		reqData      []*parseData.ReqInfo
		distDataList []*parseData.DistInfo
		errData      []*parseData.ErrInfo
		WorkerPool   = goroutine.Default()
	)
	defer c.Close()
	if err := c.CreateClient(master.OuternetIp, int(master.Port), master.RootAccount, master.RootPassword); err != nil {
		logx.Errorf("[GetReport] createSSHClient to get report err, server is %v, err is %v", master.Id, err)
		return err
	}
	//找不到报告，设置为获取报告失败
	findLocustOut, _ := c.Run("ls /data/load/master/locust_error.csv")
	if strings.TrimSpace(findLocustOut) != "/data/load/master/locust_error.csv" {
		result = common.GameGetReport
	}
	err := mr.Finish(func() (err error) {
		reqData = parseData.ReqData(c.GetDistribution(id))
		return
	}, func() (err error) {
		distDataList = parseData.DistData(c.GetRequests(id))
		return
	}, func() (err error) {
		errData = parseData.ErrData(c.GetErr())
		return
	})
	if err != nil {
		logx.Errorf("[GetReport] parseData  err, server is %v, err is %v", master.Id, err)
		return err
	}

	testNameTpsMap, err := aboutTps.AboutTps(reportID, startTime, endTime, l.svcCtx.Config.PushGateway.MonitorUrl)
	if err == nil {
		if len(testNameTpsMap) > 0 {
			var maxDistInfo parseData.DistInfo
			for _, distData = range distDataList {
				if distData.Name != "Total" {
					var tps float64
					if mapTps, ok := testNameTpsMap[distData.Name]; ok {
						tps = mapTps
					} else {
						tps = distData.Tps
					}
					totalTps += tps
					srcDistData = append(srcDistData, &parseData.DistInfo{
						ProtoType:   distData.ProtoType,
						Name:        distData.Name,
						Requests:    distData.Requests,
						Failures:    distData.Failures,
						MedianTime:  distData.MedianTime,
						AverageTime: distData.AverageTime,
						MinTime:     distData.MinTime,
						MaxTime:     distData.MaxTime,
						BodySize:    distData.BodySize,
						Tps:         tps,
					})
				} else {
					maxDistInfo.ProtoType = distData.ProtoType
					maxDistInfo.Name = distData.Name
					maxDistInfo.Requests = distData.Requests
					maxDistInfo.Failures = distData.Failures
					maxDistInfo.MedianTime = distData.MedianTime
					maxDistInfo.AverageTime = distData.AverageTime
					maxDistInfo.MinTime = distData.MinTime
					maxDistInfo.MaxTime = distData.MaxTime
					maxDistInfo.BodySize = distData.BodySize
					maxDistInfo.Tps = maxTps
				}
			}
			if totalTps, err = strconv.ParseFloat(fmt.Sprintf("%.2f", totalTps), 64); err != nil {
				logx.Errorf("failed to strconv.ParseFloat(fmt.Sprintf(\"%.2f\", totalTps, 64)"+
					"totalTps is %v err is %v", totalTps, err)
			}
			if totalTps > maxTps {
				maxDistInfo.Tps = totalTps
			}
			srcDistData = append(srcDistData, &maxDistInfo)
		}
	}
	report, _ := l.svcCtx.ReportModel.FindOne(l.ctx, reportID)
	distDataStr, _ := json.Marshal(srcDistData)
	report.DistData = sql.NullString{String: string(distDataStr)}
	reqDataStr, _ := json.Marshal(reqData)
	report.ReqData = sql.NullString{String: string(reqDataStr)}
	errDataStr, _ := json.Marshal(errData)
	report.ErrData = sql.NullString{String: string(errDataStr)}
	report.Result = int64(int(result))
	report.EndTime = sql.NullTime{Time: endTime}
	report.StartTime = sql.NullTime{Time: startTime}
	//report.MaxTps = int64(totalTps)
	err = l.svcCtx.ReportModel.Update(ctx, report)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("l.svcCtx.ReportModel.Update is fail"), "l.svcCtx.ReportModel.Update is fail err is : %+v", err)

	}
	WorkerPool.Submit(func() { //生成监控截图
		defer func() {
			recover()
		}()
		err = parseData.SaveReportPdf(l.ctx, l.svcCtx.ReportModel, reportID, common.GameReport, l.svcCtx.Config.ReportToPdfPath.Path, l.svcCtx.Config.ReportToPdfPath.UrlPre)
		if err != nil {

		}
	})
	//todo 改变机器状态
	//if err = svcCtx.MachineModel.UpdateWorkingFlag(serverIDs, common.SERVER_STATUS_IDLE); err != nil {
	//	logx.Errorf("[GetReport] update machine status err, machineIds is %v, err is %v", serverIDs, err)
	//	return err
	//}

	//ssh.KillOldTest()
	//TODO 暂不删除报告
	//c.Run("cd /data/load/master;rm -rf " + id + "_requests.csv")
	//c.Run("cd /data/load/master;rm -rf " + id + "_requests.csv")
	//c.Run("cd /data/load/master;rm -rf  " + id + "_distribution.csv")
	c.Run("rm -rf /data/load/master/locust_error.csv")
	return nil
}

func SafeRun(f func()) {
	go func() {
		if err := recover(); err != nil {
			logx.Errorf("[PanicEvent] %v", err)
		}
		f()
	}()
}
