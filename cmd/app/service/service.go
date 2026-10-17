package service

import (
	"fmt"
	"os"
	"runtime"

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
	cfg := service.Config{
		Name:        pkg.AppName,
		DisplayName: fmt.Sprintf("文件管理系统 %s", pkg.AppVersion),
		Description: "a file manager system",
	}
	// Dependencies 的语义按平台不同：
	//   - Linux(systemd)：每一项是 [Unit] 指令行，形如 "After=..." / "Requires=..."；
	//   - Windows：每一项会被 kardianos 当作「依赖的服务名」传给 SCM。
	// 之前无条件设置 systemd 指令，Windows 下服务能安装成功，但启动时 SCM 会去找名为
	// "After=local-fs.target network.target" 的服务，找不到就返回
	// "The dependency service does not exist or has been marked for deletion."。
	if runtime.GOOS == "linux" {
		cfg.Dependencies = []string{"After=local-fs.target network.target", "Requires=local-fs.target"}
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
		return boot(cfg, service)
	}
	return fmt.Errorf("目录不存在 %v", cfg.ServePath)
}

func (m *GFSService) GetAny(s string) ([]byte, []string) {
	//z.L().Sugar().Debugln("GFSService service GetAny")
	return input(), nil
}
