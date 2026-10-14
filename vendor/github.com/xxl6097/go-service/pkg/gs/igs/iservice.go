package igs

import (
	"context"

	"github.com/kardianos/service"
)

type Service interface {
	Upgrade(context.Context, string) error
	UpgradeByBuffer([]byte) error
	ClearTemp() error
	ClearAppData() error
	UnInstall() error
	RunCMD(...string) error
	Restart() error
	Status() string
}
type IService interface {
	OnConfig() *service.Config
	OnVersion() string
	OnRun(Service) error
	OnStop()
	OnShutdown()
	GetAny(string) ([]byte, []string) //数据；运行参数
	OnFinish()
	UnInstall()
}

type Installer interface {
	IService
	//OnInstall arg1: 当前运行bin文件路径，arg2安装的目标bin文件路径
	OnInstall(string, string) error
}

type DefaultService interface {
	Service
	Run() error
}
