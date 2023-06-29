package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql struct { // 数据库配置，除mysql外，可能还有mongo等其他数据库
		DataSource string // mysql链接地址，满足 $user:$password@tcp($ip:$port)/$db?$queries 格式即可
	}
	CacheRedis cache.CacheConf // redis缓存

	Git Git

	Ssh Ssh

	PushGateway     PushGateway
	ReportToPdfPath ReportToPdfPath
	Flow            Flow
}

type Git struct {
	GitCloneUrlHead string
	GitProjectPath  string
	GitlabToken     string
}

type Ssh struct {
	User           string
	MasterRemotely string
	SlaveRemotely  string
}

type PushGateway struct {
	UrlPre     string
	MonitorUrl string
}

type ReportToPdfPath struct {
	UrlPre string
	Path   string
}

type Flow struct {
	Url string
}
