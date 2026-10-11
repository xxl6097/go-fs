package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xxl6097/gofs"
	"github.com/xxl6097/gofs/githubupgrade"
)

func main() {
	cfg := gofs.Default()
	cfg.Port = 8089
	user := "admin"
	pass := "admin"
	path := "."
	cfg.ServePath = path
	// 规则里 @ 后面是 URL 路径（相对服务根目录），不是磁盘路径。
	// 填 path 等于只授权 URL /path，根目录反而无权限，列目录会 401。
	cfg.AuthRules = []string{fmt.Sprintf("%s:%s@/:rw", user, pass)}
	cfg.AllowAll = true

	up, err := githubupgrade.New(githubupgrade.Options{
		Owner:       "xxl6097",
		Repo:        "go-fs",
		AssetPrefix: "go-fs",
		DryRun:      false,
	})
	if err != nil {
		log.Fatalf("构造升级器失败: %v", err)
	}

	// 信号处理留在命令行入口：库不该替调用方决定什么时候退出。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := gofs.RunCfg(ctx, cfg, up, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, gofs.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "gofs: %v\n", err)
		os.Exit(1)
	}
}
