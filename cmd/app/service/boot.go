package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	cfg.AuthRules = []string{fmt.Sprintf("%s:%s@%s:rw", user, pass, path)}
	cfg.AllowAll = true
	bb, e := ukey.StructToGob(cfg)
	if e != nil {
		return nil
	}
	gofs.Default()
	return bb
}

func boot(cfg *gofs.Config) error {
	fmt.Printf("===>%+v\n", cfg)
	if err := runServer(cfg); err != nil {
		z.L().Sugar().Errorf("agent 退出: %v", err)
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
