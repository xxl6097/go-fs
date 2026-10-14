package core

import (
	"fmt"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"

	"path/filepath"
	"strings"
)

// Install 原始的安装文件未被删除
func Install(g igs.IService, binPath, installBinPath string) (error, []string) {
	var runArgs []string
	if gs, ok := g.(igs.Installer); ok {
		return gs.OnInstall(binPath, installBinPath), runArgs
	} else {
		cfg, args := getAny(filepath.Dir(installBinPath), g)
		if cfg != nil {
			runArgs = args
			z.L().Debug("SignFileBySelfKey", zap.String("binPath", binPath), zap.Any("cfg", cfg))
			newFilePath, e := ukey.SignFileBySelfKey(cfg, binPath)
			if e != nil {
				return e, runArgs
			}
			z.L().Debug("getAny func newFilePath", zap.String("newFilePath", newFilePath))
			binPath = newFilePath
			//return manualInstall(newFilePath, installBinPath)
		}
	}
	return manualInstall(binPath, installBinPath), runArgs
}

// 手动安装
func manualInstall(binPath, installBinPath string) error {
	if binPath == "" || installBinPath == "" {
		return fmt.Errorf("安装文件空 binPath:%v installBinPath:%v", binPath, installBinPath)
	}
	if strings.Compare(strings.ToLower(binPath), strings.ToLower(installBinPath)) == 0 {
		return fmt.Errorf("当前文件与安装文件路径一致，不允许安装 binPath:%v installBinPath:%v", binPath, installBinPath)
	} else {
		defer func() {
			_ = utils.DeleteAllDirector(binPath)
		}()
	}
	err := utils.Copy(binPath, installBinPath)
	if err != nil {
		z.L().Error("文件拷贝失败", zap.Error(err))
		return err
	}
	return nil
}

func getAny(binDir string, g igs.IService) ([]byte, []string) {
	if g == nil {
		z.L().Error("igs.IService is nil")
		return nil, nil
	}
	buffer, runArgs := g.GetAny(binDir)
	if buffer == nil {
		cfg := map[string]any{"path": binDir}
		bb, err := ukey.StructToGob(cfg)
		if err != nil {
			return []byte(err.Error()), runArgs
		}
		buffer = bb
	}
	return buffer, runArgs
}

//func getAny(binDir string, g igs.IService) any {
//	if g == nil {
//		glog.Error("igs.IService is nil")
//		return nil
//	}
//	cfg := g.GetAny(filepath.Dir(binDir))
//	if cfg == nil {
//		cfg = ukey.KeyBuffer{}
//	}
//	switch v := cfg.(type) {
//	case ukey.KeyBuffer:
//		v.MenuDisable = true
//		return v
//	case *ukey.KeyBuffer:
//		v.MenuDisable = true
//		return v
//	}
//	return cfg
//}

//安装、升级、卸载、开启、关闭、重启分别定义接口并继承BaseService
//1、安装，无实现，执行默认行为1（复制）；有实现，且参数合法，则执行子类行为，否则执行默认行为1
//2、作为go-service，默认行为2需要做签名功能
