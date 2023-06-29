/*
@Time : 2023/4/20 18:14
@Author : Hhx06
@File : sftp
@Description: sftp
@Software: GoLand
*/

package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	ssh2 "single/common/ssh"
	"time"
)

type FTPClient struct {
	Host       string       `json:"host"`
	User       string       `json:"user"`
	Pwd        string       `json:"pwd"`
	sshClient  *ssh.Client  //ssh client
	SftpClient *sftp.Client //sftp client
	LastResult string       //最近一次运行的结果
	Port       int          `json:"port"`
}

func (cliConf *FTPClient) CreateClient(host string, port int, username, password string) error {
	var (
		sshClient  *ssh.Client
		SftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.User = username
	cliConf.Pwd = password
	cliConf.Port = port

	config := ssh.ClientConfig{
		User:            cliConf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	publicKeyAuthFunc, err := ssh2.PublicKeyAuthFunc("")
	if err == nil {
		config.Auth = append(config.Auth, publicKeyAuthFunc)
	}

	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)
	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		return err
	}
	cliConf.sshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建SftpClient
	if SftpClient, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	cliConf.SftpClient = SftpClient
	return nil
}

func (cliConf *FTPClient) Run(shell string) (string, error) {
	var (
		session *ssh.Session
		err     error
	)
	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.sshClient.NewSession(); err != nil {
		return "", err
	}
	defer session.Close()
	//执行shell
	if output, err := session.CombinedOutput(shell); err != nil {
		return string(output), err
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult, nil
}

func (cliConf *FTPClient) Upload(srcPath, dstPath string) error {
	srcFile, _ := os.Open(srcPath) //本地
	_ = cliConf.SftpClient.Remove(dstPath)
	dstFile, _ := cliConf.SftpClient.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1048576)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error occurred:", err)
				return err
			} else {
				break
			}
		}
		_, err = dstFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (cliConf *FTPClient) Download(srcPath, dstPath string) {
	cliConf.CreateClient(cliConf.Host, cliConf.Port, cliConf.User, cliConf.Pwd)
	srcFile, err := cliConf.SftpClient.Open(srcPath) //远程
	if err != nil {
		fmt.Println(err)
		log.Fatalf("创建客户端失败")
		return
	}
	dstFile, _ := os.Create(dstPath) //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Fatalln("error occurred", err)
	}
	fmt.Println("文件下载完毕")
}

func (cliConf *FTPClient) GetDistribution(id string) string {
	dis, err := cliConf.Run("cd /data/load/master; cat " + id + "_distribution.csv")
	if err != nil {
		logx.Errorf("[getReport.GetDistribution] err is %v, msg is %v", err, dis)
		return ""
	}

	return dis
}
func (cliConf *FTPClient) GetRequests(id string) string {
	req, err := cliConf.Run("cd /data/load/master; cat " + id + "_requests.csv")
	if err != nil {
		logx.Errorf("[getReport.GetRequests] err is %v, msg is %v", err, req)
		return ""
	}

	return req
}

func (cliConf *FTPClient) GetErr() string {
	e, err := cliConf.Run("cd /data/load/master; cat  locust_error.csv")
	if err != nil {
		logx.Errorf("[getReport.GetErrData] err is %v, msg is %v", err, e)
		return ""
	}
	return e
}

func (cliConf *FTPClient) Close() {
	var (
		err error
	)

	err = cliConf.SftpClient.Close()
	if err != nil {
		logx.Errorf("sftpClient Close Err, err is %v", err)
	}
	err = cliConf.sshClient.Close()
	if err != nil {
		logx.Errorf("sshClient Close Err, err is %v", err)
	}
}
