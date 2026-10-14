package internal

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"

	"os"
)

func init() {
	InitLog()
}

func (c *CoreService) Start(s service.Service) error {
	status, err := s.Status()
	//defer glog.Flush()
	z.L().Info(fmt.Sprintf("启动服务【%s】", s.String()), zap.Any("status", status), zap.Any("Platform", s.Platform()), zap.Error(err))
	go func() {
		if c.isrv != nil {
			e := c.isrv.OnRun(c)
			if e != nil {
				z.L().Warn("运行失败", zap.Error(e))
			} else {
				z.L().Info("运行成功")
			}
		}
	}()
	return nil
}

func (c *CoreService) Stop(s service.Service) error {
	//defer glog.Flush()
	status, err := s.Status()
	ok := service.Interactive()
	//注意，这个地方在非windows下不行！
	//c.clearForUninstall()
	z.L().Info(fmt.Sprintf("停止服务【%s】", s.String()), zap.Any("status", status), zap.Any("Platform", s.Platform()), zap.Error(err))

	if c.isrv != nil {
		c.isrv.OnStop()
	}
	if ok {
		z.L().Warn("停止deamon")
		os.Exit(0)
	}
	return nil
}

func (c *CoreService) Shutdown(s service.Service) error {
	//defer glog.Flush()
	if c.isrv != nil {
		c.isrv.OnShutdown()
	}
	status, err := s.Status()
	z.L().Info(fmt.Sprintf("Shutdown【%s】", s.String()), zap.Any("status", status), zap.Any("Platform", s.Platform()), zap.Error(err))
	return nil
}
