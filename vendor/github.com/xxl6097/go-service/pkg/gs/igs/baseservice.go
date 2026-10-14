package igs

import "github.com/xxl6097/glog/pkg/z"

type BaseService struct {
}

func (q *BaseService) OnStop() {
	z.L().Sugar().Debugln("service OnStop")
}

func (q *BaseService) OnShutdown() {
	z.L().Sugar().Debugln("service OnShutdown")
}
func (q *BaseService) OnFinish() {
	z.L().Sugar().Debugln("service OnFinish")
}
func (q *BaseService) UnInstall() {
	z.L().Sugar().Debugln("service UnInstall")
}
