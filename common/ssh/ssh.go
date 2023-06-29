/*
@Time : 2023/4/20 17:06
@Author : Hhx06
@File : ssh
@Description: 操作ssh
@Software: GoLand
*/

package ssh

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"single/stressTask/model"
	"time"
)

type Service struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Pwd     string `json:"pwd"`
	Type    string `json:"type"`
	KeyPath string `json:"key_path"`
	Port    int64  `json:"port"`
	Client  *ssh.Client
}

func (s *Service) Connect() (session *ssh.Session, err error) {
	var authFunc ssh.AuthMethod

	config := &ssh.ClientConfig{
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second,
	}

	authFunc, err = PublicKeyAuthFunc("")
	if err == nil {
		config.Auth = append(config.Auth, authFunc)
	}

	s.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), config)
	if err != nil {
		return nil, err
	}

	session, err = s.Client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func PublicKeyAuthFunc(kPath string) (ssh.AuthMethod, error) {
	if len(kPath) == 0 {
		kPath = "~/.ssh/id_rsa"
	}
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		logx.Error("find key's home dir failed", err)
		return nil, err
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logx.Error("ssh key file read failed", err)
		return nil, err
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logx.Error("ssh key signer failed", err)
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// KillOldTest 清理locust进程
func KillOldTest(masterServer *model.Machine, slaveServer []*model.Machine) error {
	var (
		killLocustCommand = `killall -9 -q locust`
		killSlaveCommand  = "pid=`ps -ef|grep /data/load/slave |grep -v grep|awk '{print $2}' | xargs`; if [ $pid ]; then  kill -9 $pid;  fi"
	)
	err := mr.Finish(func() error {
		if masterServer != nil { //kill master server
			masterSSHServer := Service{Host: masterServer.OuternetIp, User: masterServer.RootAccount, Pwd: masterServer.RootPassword, Type: "password", KeyPath: "", Port: masterServer.Port}
			session, err := masterSSHServer.Connect()
			if err != nil || session == nil {
				logx.Errorf("[KillOldTest]  master connect is  err: %v", err)
				return err
			}
			combinedOutput, err := session.CombinedOutput(killLocustCommand)
			if err != nil && string(combinedOutput) != "" {
				logx.Errorf("killLocustCommand run combinedOutput: %s, err %v", string(combinedOutput), err)
				return errors.New("locust 进程kill失败")
			}
			session.Close()
			masterSSHServer.Client.Close()
		}
		return nil
	}, func() (err error) {
		if len(slaveServer) != 0 { // kill slave
			for _, machine := range slaveServer { // 连接 slave
				slaveSSHServer := Service{Host: machine.OuternetIp, User: machine.RootAccount, Pwd: machine.RootPassword, Type: "password", KeyPath: "", Port: machine.Port}
				session, err := slaveSSHServer.Connect()
				if err != nil || session == nil {
					logx.Errorf("slave connect is  err: %v", err)
					return err
				}

				combinedOutput, err := session.CombinedOutput(killSlaveCommand)
				if err != nil && string(combinedOutput) != "" {
					logx.Errorf("killLocustCommand run combinedOutput: %s, err %v", string(combinedOutput), err)
					return errors.New("locust 进程kill失败")
				}
				session.Close()
				slaveSSHServer.Client.Close()
			}
		}
		return nil
	})

	if err != nil {
		return nil
	}

	return nil
}
