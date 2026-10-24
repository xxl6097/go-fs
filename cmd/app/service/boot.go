package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/xxl6097/gfs/internal"
	"github.com/xxl6097/gfs/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/gofs"
)

func input() []byte {
	cfg := gofs.Default()
	cfg.Bind = utils.InputStringEmpty("服务监听地址(0.0.0.0):", "0.0.0.0")
	cfg.Port = utils.InputIntDefault("管理后台端口(8088):", 8088)
	user := utils.InputStringEmpty("管理后台用户名(admin):", "admin")
	pass := utils.InputString("请输入管理员密码：")
	path, err := InputPath("请输入文件管理路径：")
	if err != nil {
		return nil
	}
	cfg.ServePath = path
	// 规则里 @ 后面是 URL 路径（相对服务根目录），不是磁盘路径。
	// 填 path 等于只授权 URL /path，根目录反而无权限，列目录会 401。
	cfg.AuthRules = []string{fmt.Sprintf("%s:%s@/:rw", user, pass)}
	cfg.AllowAll = true
	bb, e := ukey.StructToGob(cfg)
	if e != nil {
		return nil
	}
	return bb
}

func boot(cfg *gofs.Config, gs igs.Service) error {
	fmt.Printf("启动参数：%+v\n", cfg)
	if err := runServer(cfg, gs); err != nil {
		z.L().Sugar().Errorf("gfs启动失败: %v", err)
		return err
	}
	return nil
}

func loadBuffer() (*gofs.Config, error) {
	byteArray, err := ukey.Load()
	if err != nil {
		return nil, err
	}
	cfg, err := decodeConfig(byteArray)
	if err != nil {
		z.L().Sugar().Error("解密错误", err)
		return nil, err
	}
	return cfg, nil
}

// decodeConfig 从默认值起解码签名里的 gob 配置。
//
// 旧版本签进二进制的配置没有后来新增的字段（例如 AuthLockMaxDelay），
// 从零值解码会让这些字段变成 0，过不了 gofs 的校验；gob 只覆盖流里
// 存在的字段，所以缺失的新字段会保留这里的默认值。
func decodeConfig(byteArray []byte) (*gofs.Config, error) {
	cfg := gofs.Default()
	if err := ukey.GobToStruct(byteArray, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func runServer(cfg *gofs.Config, gs igs.Service) error {
	cfg.Version = pkg.AppVersion
	if err := internal.RunCfg(context.Background(), cfg, gs, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, gofs.ErrHelp) {
			os.Exit(0)
		}
		z.L().Sugar().Errorf("gofs: %v\n", err)
		return err
	}
	return nil
}

func checkDir(path string) bool {
	for {
		mounted, err := DirExists(path)
		if err == nil && mounted {
			return true
		}
		time.Sleep(3 * time.Second)
	}
}

// DirExists 判断路径是否存在，并且是目录（支持软链接）
func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}
