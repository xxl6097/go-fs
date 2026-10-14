package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/internal/core"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/go-service/pkg/utils/util"
	"go.uber.org/zap"
)

func (this *CoreService) install() error {
	//this.reqeustWindowsUser()
	s, err := this.statusService()
	var isRemoved = true
	if err == nil {
		no := utils.InputString(fmt.Sprintf("%s%s%s", "检测到", this.config.Name, "程序已经安装，卸载/更新/取消?(y/u/n):"))
		switch no {
		case "y", "Y", "Yes", "YES":
			isRemoved = true
			e := this.uninstall()
			if e != nil {
				z.L().Warn("卸载失败", zap.Any("config", this.config), zap.Error(e))
			} else {
				z.L().Info("卸载成功", zap.Any("config", this.config))
			}
			break
		case "u", "U", "Update", "UPDATE":
			isRemoved = false
			e := this.stopService()
			if err != nil {
				//return err
				z.L().Error("升级", zap.Error(e))
			}
			break
		default:
			z.L().Debug("结束安装")
			time.Sleep(time.Second * 3)
			os.Exit(0)
			return err
		}
	} else if s != service.StatusUnknown {
		e := this.uninstallService()
		if e != nil {
			z.L().Warn("卸载失败", zap.Error(e))
		}
	}
	util.SetFirewall(this.config.Name, this.config.Executable)
	e1 := util.SetRLimit()
	if e1 != nil {
		z.L().Warn("SetRLimit", zap.Error(e1))
	}

	if isRemoved {
		_ = utils.ResetDirector(this.workDir)
	} else {
		_ = utils.CheckDirector(this.workDir)
	}

	//这个地方是取的当前运行的执行文件
	currentBinPath, e := os.Executable()
	if e != nil {
		z.L().Fatal("s.Executable() error", zap.Error(e))
		return e
	}

	ee, runArgs := core.Install(this.isrv, currentBinPath, this.config.Executable)
	if ee != nil {
		z.L().Error("Install err", zap.Error(ee))
		return ee
	}

	e = os.Chdir(this.workDir)
	if e != nil {
		z.L().Error("Chdir err", zap.Error(e))
		return e
	}

	err = this.installService(runArgs)
	if err == nil {
		z.L().Debug(fmt.Sprintf("[%s]安装成功!", this.config.DisplayName))
	} else {
		z.L().Warn(fmt.Sprintf("[%s]安装失败", this.config.DisplayName), zap.Error(err))
	}
	time.Sleep(time.Second * 1)
	err = this.startService()
	if err != nil {
		z.L().Warn(fmt.Sprintf("[%s]启动失败，错误信息:%v\n", this.config.DisplayName), zap.Error(err))
	} else {
		z.L().Debug(fmt.Sprintf("[%s]启动成功!", this.config.DisplayName))
	}
	time.Sleep(time.Second * 1)
	z.L().Info(this.Status())
	err = os.RemoveAll(currentBinPath)
	if err != nil {
		fmt.Println("delete error", currentBinPath, err)
	}
	return nil
}
func (this *CoreService) uninstall() error {
	defer func() {
		if this.isrv != nil {
			if gs, ok := this.isrv.(igs.Installer); ok {
				gs.UnInstall()
			}
		}
		this.clearForUninstall()
		z.L().Debug("尝试停止服务")
		err := this.stopService()
		z.L().Debug("2尝试停止服务", zap.Error(err))
		this.clearForUninstall()
	}()

	if utils.IsOpenWRT() {
		err := this.stopService()
		if err != nil {
			z.L().Warn(fmt.Sprintf("[%s]停止失败", this.config.Name), zap.Error(err))
		} else {
			z.L().Info(fmt.Sprintf("[%s]停止成功", this.config.Name))
		}
	}
	err := this.uninstallService()
	if err != nil {
		z.L().Warn(fmt.Sprintf("[%s]卸载失败 ", this.config.Name), zap.Error(err))
	} else {
		z.L().Info(fmt.Sprintf("[%s]成功卸载", this.config.Name))
	}
	return err
}

func (this *CoreService) patchUpgrade(ctx context.Context, binUrlOrLocal string) error {
	downFilePath, err := utils.CheckFileOrDownload(ctx, binUrlOrLocal)
	if err != nil {
		z.L().Warn("升级失败", zap.Error(err))
		return err
	}

	if !utils.FileExists(downFilePath) {
		return fmt.Errorf("差分升级文件不存在 %s", downFilePath)
	}
	return this.update(downFilePath, true)
}

// 1. 检测升级文件是本地还是网络文件（下载）；
// 2. 升级文件最终需要被删除，所以使用defer删除；
// 3. 给升级文件赋予0755权限
func (this *CoreService) upgrade(ctx context.Context, binUrlOrLocal string) error {
	downFilePath, err := utils.CheckFileOrDownload(ctx, binUrlOrLocal)
	if err != nil {
		z.L().Warn("升级失败", zap.Error(err))
		return err
	}
	z.L().Sugar().Debug("升级文件 ", downFilePath)
	var signFilePath string
	patch := false
	if filepath.Ext(downFilePath) == ".patch" {
		signFilePath = downFilePath
		patch = true
		z.L().Sugar().Debug("差量升级文件 ", downFilePath)
	} else {
		z.L().Sugar().Debug("全量升级文件，签名  ", downFilePath)
		tempFilePath, e := ukey.SignFileByOldFileKey(this.config.Executable, downFilePath)
		//签名完后会生产出新的签名文件，那么下载的文件需要被删除
		_ = utils.DeleteAllDirector(downFilePath)
		if e != nil {
			z.L().Error("升级签名错误", zap.Error(e))
			return err
		}
		signFilePath = tempFilePath
	}

	if !utils.FileExists(signFilePath) {
		return fmt.Errorf("升级文件不存在 %s", signFilePath)
	}

	return this.update(signFilePath, patch)
	//err = os.Chmod(signFilePath, 0755)
	//if err != nil {
	//	eMsg := fmt.Sprintf("赋权限错误: %s %v\n", signFilePath, err)
	//	glog.Error(eMsg)
	//	return fmt.Errorf(eMsg)
	//}
	//glog.Println("当前进程ID:", os.Getpid(), this.config.Executable)
	//err = utils.PerformUpdate(signFilePath, this.config.Executable, patch)
	////同样，更新完后，需要删除签名文件
	//_ = utils.DeleteAllDirector(filepath.Dir(filepath.Dir(signFilePath)))
	//if err != nil {
	//	glog.Error("升级失败", err)
	//	return err
	//}
	//glog.Error("升级成功")
	//if utils.IsWindows() {
	//	glog.Debug(utils.RunCmd("dir"))
	//} else {
	//	glog.Debug(utils.RunCmd("ls", "-l"))
	//}
	//return this.RunCMD("restart")
}

func (this *CoreService) changeSelf(buffer []byte) error {
	if buffer == nil || len(buffer) == 0 {
		return errors.New("配置buffer空")
	}
	binFilePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前可执行文件路径出错: %v\n", err)
	}

	tempFilePath, e := ukey.SignFileByBuffer(buffer, binFilePath)
	if e != nil {
		z.L().Warn("升级签名错误", zap.Error(e))
		return err
	}
	signFilePath := tempFilePath
	if !utils.FileExists(signFilePath) {
		return fmt.Errorf("自升级文件不存在 %s", signFilePath)
	}
	return this.update(signFilePath, false)
}

func (this *CoreService) update(signFilePath string, patch bool) error {
	err := os.Chmod(signFilePath, 0755)
	if err != nil {
		eMsg := fmt.Sprintf("赋权限错误: %s %v\n", signFilePath, err)
		z.L().Warn(eMsg)
		return fmt.Errorf(eMsg)
	}
	upgradeName := "全量更新"
	if !patch {
		err = utils.IsMissMatchOsApp(signFilePath)
		if err != nil {
			return err
		}
	} else {
		upgradeName = "差量更新"
	}
	z.L().Debug("当前进程", zap.Int("pid", os.Getpid()), zap.String("upgradeName", upgradeName), zap.String("executable", this.config.Executable))
	err = utils.PerformUpdate(signFilePath, this.config.Executable, patch)
	//同样，更新完后，需要删除签名文件
	defer func() {
		_ = utils.DeleteAllDirector(signFilePath)
	}()
	if err != nil {
		z.L().Error("升级失败", zap.String("upgradeName", upgradeName), zap.Error(err))
		//_ = utils.CopyToTemp(signFilePath)
		return err
	}
	z.L().Debug("升级成功", zap.String("upgradeName", upgradeName))
	if utils.IsWindows() {
		z.L().Debug(utils.RunCmd("dir"))
	} else {
		z.L().Debug(utils.RunCmd("ls", "-l"))
	}
	// 只拉起重启子进程，不要在这里 stopService()：当前服务自行 systemctl stop
	// （KillMode=control-group）会把同一 cgroup 里的 systemd-run 辅助进程、以及
	// systemctl 客户端自身一并 SIGTERM，restart 子进程还没被拉起服务就停在停止态了。
	// restart 子进程在独立 scope 里执行 systemctl restart（stop+start），
	// 当前进程会由那次 SIGTERM 触发优雅退出。
	err = this.RunCMD("restart")
	return err
}

func (this *CoreService) clearForUninstall() {
	_ = utils.DeleteAllDirector(this.workDir)
	appDir := zutil.AppHome()
	_ = utils.DeleteAllDirector(appDir)
}
func (this *CoreService) clearAppData() error {
	appDir := zutil.AppHome()
	return utils.DeleteAllDirector(appDir)
}
