package service

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/xxl6097/gfs/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/gofs"
)

type GFSService struct {
	igs.BaseService
	gs igs.Service
}

func (m *GFSService) UnInstall() {
	z.L().Sugar().Debugln("mu service OnUninstall")
	keyDir := gofs.DefaultKeyDir()
	if keyDir != "" {
		_ = os.RemoveAll(keyDir)
	}
}

func (m *GFSService) OnConfig() *service.Config {
	z.L().Sugar().Debugln("mu service OnConfig")
	cfg := service.Config{
		Name:         pkg.AppName,
		DisplayName:  fmt.Sprintf("文件管理系统 %s", pkg.AppVersion),
		Description:  "a file manager system",
		Dependencies: []string{"After=local-fs.target network.target", "Requires=local-fs.target"},
	}
	// macOS 下安装目录在用户目录（见 util_darwin.go），服务也注册为用户级
	// LaunchAgent（~/Library/LaunchAgents/<name>.plist）；
	// 否则 kardianos 会写到 /Library/LaunchDaemons，仍然需要 root。
	if utils.IsMacOs() {
		cfg.Option = service.KeyValue{"UserService": true}
	}
	return &cfg
}

func (m *GFSService) OnVersion() string {
	z.L().Sugar().Debugln("mu service OnVersion")
	pkg.Version()
	return pkg.AppVersion
}

func (m *GFSService) OnRun(service igs.Service) error {
	z.L().Sugar().Debugln("GFSService service OnRun")
	m.gs = service
	cfg, err := loadBuffer()
	if err != nil {
		return err
	}
	if checkDir(cfg.ServePath) {
		return boot(cfg)
	}
	return fmt.Errorf("目录不存在 %v", cfg.ServePath)
}

func (m *GFSService) GetAny(s string) ([]byte, []string) {
	z.L().Sugar().Debugln("GFSService service GetAny")
	return input(), nil
}
