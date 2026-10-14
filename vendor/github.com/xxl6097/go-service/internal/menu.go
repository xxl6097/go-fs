package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"

	"os"
)

func (this *CoreService) menu() error {
	if len(os.Args) > 1 {
		k := os.Args[1]
		err := this.runSwitch(k)
		if !errors.Is(err, errors.ErrUnsupported) {
			return err
		}
	}
	if ukey.CanShowMenu() || !this.isServiceApp() {
		//glog.Debug("运行菜单", os.Getpid())
		return this.runMenu()
	}
	z.L().Debug("运行参数", zap.Strings("Args", os.Args), zap.Int("pid", os.Getpid()))
	return this.runService()
}

func (this *CoreService) runSwitch(cmd string) error {
	switch cmd {
	case "-v", "v", "version":
		fmt.Println(this.isrv.OnVersion())
		return nil
	case "install":
		return this.install()
	case "uninstall":
		return this.uninstall()
	case "upgrade", "update":
		if len(os.Args) <= 2 {
			return fmt.Errorf("参数错误，缺少升级文件链接 %v", os.Args)
		}
		binUrlOrLocal := os.Args[2]
		if binUrlOrLocal == "" {
			return errors.New("文件信息空～")
		}
		return this.upgrade(context.Background(), binUrlOrLocal)
	case "start":
		return this.startService()
	case "stop":
		return this.stopService()
	case "restart":
		return this.restartService()
	case "status":
		z.L().Debug(this.Status())
		return nil
	case "r", "run":
		return this.isrv.OnRun(this)
	}
	return errors.ErrUnsupported
}

func (this *CoreService) runMenu() error {
	defer func() {
		if this.isrv != nil {
			this.isrv.OnFinish()
		}
		if utils.IsWindows() {
			utils.ExitAnyKey()
		} else {
			utils.ExitCountDown(1)
		}
	}()
	keys := []string{"install", "uninstall", "upgrade", "restart", "stop", "status", "v"}
	fmt.Println("1. 安装程序")
	fmt.Println("2. 卸载程序")
	fmt.Println("3. 升级程序")
	fmt.Println("4. 重启程序")
	fmt.Println("5. 停止程序")
	fmt.Println("6. 服务状态")
	fmt.Println("7. 查看版本")
	index := utils.InputInt("请根据菜单选择：")
	if index >= 1 && index <= len(keys) {
		index--
		key := keys[index]
		return this.runSwitch(key)
	} else {
		fmt.Println("未知选项", index)
	}
	return nil
}

//func (this *CoreService) runMenu1() error {
//	defer func() {
//		utils.Exit()
//	}()
//	fmt.Println("1. 安装程序")
//	fmt.Println("2. 卸载程序")
//	fmt.Println("3. 升级程序")
//	fmt.Println("4. 重启程序")
//	fmt.Println("5. 停止程序")
//	fmt.Println("6. 查看版本")
//	index := utils.InputInt("请根据菜单选择：")
//	switch index {
//	case 1:
//		return this.install()
//	case 2:
//		return this.uninstall()
//	case 3:
//		if len(os.Args) <= 2 {
//			return fmt.Errorf("参数错误，缺少升级文件链接 %v", os.Args)
//		}
//		binUrlOrLocal := os.Args[2]
//		if binUrlOrLocal == "" {
//			return errors.New("文件信息空～")
//		}
//		return this.upgrade(context.Background(), binUrlOrLocal)
//	case 4:
//		return this.restartService()
//	case 5:
//		return this.stopService()
//	case 6:
//		fmt.Println(this.isrv.OnVersion())
//		break
//	default:
//		fmt.Println("未知选项", index)
//		break
//	}
//	return nil
//}

//func (this *CoreService) menu() error {
//	if len(os.Args) > 1 {
//		k := os.Args[1]
//		switch k {
//		case "-v", "v", "version":
//			fmt.Println(this.isrv.OnVersion())
//			return nil
//		case "install":
//			return this.install()
//		case "uninstall":
//			return this.uninstall()
//		case "upgrade", "update":
//			if len(os.Args) <= 2 {
//				return fmt.Errorf("参数错误，缺少升级文件链接 %v", os.Args)
//			}
//			binUrlOrLocal := os.Args[2]
//			if binUrlOrLocal == "" {
//				return errors.New("文件信息空～")
//			}
//			return this.upgrade(context.Background(), binUrlOrLocal)
//		case "start":
//			return this.startService()
//		case "stop":
//			return this.stopService()
//		case "restart":
//			return this.restartService()
//		case "r", "run":
//			return this.isrv.OnRun(this)
//		}
//	}
//	if ukey.CanShowMenu() {
//		glog.Debug("运行菜单", os.Getpid())
//		return this.runMenu()
//	}
//	return this.runService()
//}
