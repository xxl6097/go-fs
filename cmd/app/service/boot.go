package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/gofs"
)

func input() []byte {
	cfg := gofs.Default()
	cfg.Port = utils.InputInt("请输入端口：")
	user := utils.InputString("请输入管理员账号：")
	pass := utils.InputString("请输入管理员密码：")
	path := utils.InputString("请输入文件管理路径：")
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

func boot(cfg *gofs.Config) error {
	fmt.Printf("启动参数：%+v\n", cfg)
	if err := runServer(cfg); err != nil {
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
	var cfg gofs.Config
	err = ukey.GobToStruct(byteArray, &cfg)
	if err != nil {
		fmt.Println("解密错误", err)
		return nil, err
	}
	return &cfg, nil
}

func runServer(cfg *gofs.Config) error {
	// 信号处理留在命令行入口：库不该替调用方决定什么时候退出。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := gofs.RunCfg(ctx, cfg, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, gofs.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "gofs: %v\n", err)
		//os.Exit(1)
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
