/*
@Time : 2023/4/20 15:28
@Author : Hhx06
@File : git
@Description: 操作git
@Software: GoLand
*/

package git

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/exec"
	"strings"
)

//CloneCommand 执行git项目拉取（远程指令） //增加分支名  projectPath = l.svcCtx.Config.Git.GitProjectPath gitCloneUrlHead = l.svcCtx.Config.Git.GitCloneUrlHead
func CloneCommand(gitPath, token, projectPath, gitCloneUrlHead string) (string, error) {
	gitPathArr := strings.Split(gitPath, "/")

	gitPath = strings.TrimLeft(strings.TrimSpace(gitPath), "https://")
	gitClonePath := fmt.Sprintf("%s%s@%s.git", gitCloneUrlHead, token, gitPath)
	gitSavePath := projectPath + gitPathArr[len(gitPathArr)-1]
	//判断是否存在，存在则删除重新克隆
	_, exist := os.Stat(gitSavePath)
	if exist == nil || os.IsExist(exist) {
		if err := os.RemoveAll(gitSavePath); err != nil {
			logx.Errorf("原数据删除失败")
			return "", err
		}
	}

	cmdGit := exec.Command("git", "clone", gitClonePath, gitSavePath)
	msg, err := cmdGit.CombinedOutput()
	return string(msg), err
}

//CheckoutCommand 切换项目目标分支 projectPath = l.svcCtx.Config.Git.GitProjectPath
func CheckoutCommand(gitPath, branch, projectPath string) (string, error) {
	gitPathArr := strings.Split(gitPath, "/")
	gitSavePath := projectPath + gitPathArr[len(gitPathArr)-1]
	if branch == "master" {
		return "", nil
	}
	cmdGit := exec.Command("git", "checkout", branch)
	cmdGit.Dir = gitSavePath
	msg, err := cmdGit.CombinedOutput()

	return string(msg), err
}

// CloneDir 操作拉取项目到指定目录 编译
func CloneDir(url, gitlabToken, branch, entryDir, projectPath, gitCloneUrlHead string) (string, error) {
	var dirPath string

	if msg, err := CloneCommand(url, gitlabToken, projectPath, gitCloneUrlHead); err != nil {
		logx.Errorf("gitCloneDir err,return Msg:%s branch: %s err:%v", msg, branch, err)
		return dirPath, errors.New("git仓库本地拉取失败")
	}

	// 编译指定test目录下的文件
	gitPathArr := strings.Split(url, "/")
	newProjectPath := projectPath + gitPathArr[len(gitPathArr)-1]

	if msg, err := CheckoutCommand(url, branch, projectPath); err != nil {
		logx.Errorf(" GitCheckoutCommand err, gitPath:%s, Branch:%s, return msg:%s, err:%v",
			url, branch, msg, err)
		return dirPath, errors.New("目标分支切换失败")
	}

	path := newProjectPath + "/test/" + entryDir
	command := `cd ` + path + ";" + ` sudo /usr/local/go/bin/go get ./... && sudo /usr/local/go/bin/go build -ldflags="-s -w"  main.go`
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Env = []string{"CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64"}
	if output, err := cmd.CombinedOutput(); err != nil {
		logx.Errorf("cmd.CombinedOutput err, gitPath:%s, return Msg:%s,err:%v", url, string(output), err)
		return dirPath, errors.New("压测任务编译失败")
	}
	dirPath = projectPath + "/test/" + entryDir + "/"
	return dirPath, nil
}
