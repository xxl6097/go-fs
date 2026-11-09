package main

import (
	"os"

	_ "github.com/xxl6097/gfs/cmd/app/buffer"
	"github.com/xxl6097/gfs/cmd/app/service"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs"
	"go.uber.org/zap"
)

func init() {
	// 在运行时初始化前禁用 Huge Pages 支持，避免 getHugePageSize() panic
	_ = os.Setenv("GODEBUG", "madvdontneed=1")
}

func main() {
	err := gs.Run(&service.GFSService{})
	z.L().Debug("程序结束", zap.Error(err))
}
