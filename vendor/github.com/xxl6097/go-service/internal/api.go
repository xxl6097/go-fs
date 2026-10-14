package internal

import (
	"context"
	"fmt"
	"github.com/kardianos/service"
	"github.com/xxl6097/go-service/pkg/utils"
)

func (this *CoreService) Upgrade(ctx context.Context, binUrl string) error {
	return this.upgrade(ctx, binUrl)
}

func (this *CoreService) PatchUpgrade(ctx context.Context, binUrl string) error {
	return this.patchUpgrade(ctx, binUrl)
}

func (this *CoreService) UpgradeByBuffer(buffer []byte) error {
	return this.changeSelf(buffer)
}

func (this *CoreService) ClearTemp() error {
	return utils.ClearTemp()
}

func (this *CoreService) ClearAppData() error {
	return this.clearAppData()
}

func (this *CoreService) UnInstall() error {
	//return this.uninstall()
	return this.RunCMD("uninstall")
}

func (this *CoreService) RunCMD(args ...string) error {
	return utils.RunCmdBySelf(this.config.Executable, args...)
}

func (this *CoreService) Restart() error {
	return this.RunCMD("restart")
}

func (this *CoreService) Status() string {
	s, e := this.statusService()
	if e != nil {
		return e.Error()
	}
	if s == service.StatusUnknown {
		return fmt.Sprintf("%s 服务未安装", this.config.Name)
	} else if s == service.StatusRunning {
		return fmt.Sprintf("%s 服务运行中...", this.config.Name)
	} else if s == service.StatusStopped {
		return fmt.Sprintf("%s 服务已停止", this.config.Name)
	} else {
		return fmt.Sprintf("%s 服务未知状态", this.config.Name)
	}
}
