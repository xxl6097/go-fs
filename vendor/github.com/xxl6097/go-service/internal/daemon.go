package internal

import (
	"errors"
	"fmt"

	"github.com/kardianos/service"
	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

func (this *CoreService) createService() error {
	s, e := service.New(this, this.config)
	if e != nil {
		return fmt.Errorf("新建服务错误:%+v", e)
	}
	if s == nil {
		return errors.New("新建服务失败～")
	}
	service.Interactive()
	this.srv = s
	return nil
}

func (this *CoreService) control(cmd string) error {
	z.L().Debug(fmt.Sprintf("[%s]正在%s", this.config.Name, cmd))
	e := service.Control(this.srv, cmd)
	if e != nil {
		z.L().Debug(fmt.Sprintf("[%s]%s失败", this.config.Name, cmd), zap.Error(e))
		return e
	}
	z.L().Debug(fmt.Sprintf("[%s]%s成功", this.config.Name, cmd))
	return nil
}

func (this *CoreService) statusService() (service.Status, error) {
	return this.srv.Status()
}
func (this *CoreService) IsRunning() bool {
	if this.srv == nil {
		return false
	}
	status, err := this.statusService()
	if err != nil {
		//glog.Println(err)
		return false
	}
	if status == service.StatusRunning {
		//glog.Println(this.config.Name, "is running")
		return true
	} else if status == service.StatusStopped {
		//glog.Println(this.config.Name, "is stopped")
	} else {
		//glog.Println(this.config.Name, "StatusUnknown", status)
	}
	return false
}

func (this *CoreService) uninstallService() error {
	return this.control(service.ControlAction[4])
}
func (this *CoreService) installService(runArgs []string) error {
	if this.config.Option == nil {
		this.config.Option = make(map[string]interface{})
	}
	//windows下，服务->登录-登录身份->本地系统账户->允许服务与桌面交互
	this.config.Option["Interactive"] = true
	if runArgs != nil && len(runArgs) > 0 {
		this.config.Arguments = append(this.config.Arguments, runArgs...)
		s, e := service.New(this, this.config)
		if e != nil {
			z.L().Error("service New err", zap.Error(e))
			return e
		}
		this.srv = s
	}
	return this.control(service.ControlAction[3])
}
func (this *CoreService) restartService() error {
	return this.control(service.ControlAction[2])
}
func (this *CoreService) stopService() error {
	return this.control(service.ControlAction[1])
}
func (this *CoreService) startService() error {
	return this.control(service.ControlAction[0])
}
func (this *CoreService) runService() error {
	return this.srv.Run()
}
