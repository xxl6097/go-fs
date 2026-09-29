package service

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/xxl6097/gfs/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
)

type GFSService struct {
	igs.BaseService
	gs igs.Service
}

func (m *GFSService) OnConfig() *service.Config {
	z.L().Sugar().Debugln("mu service OnConfig")
	cfg := service.Config{
		Name:        pkg.AppName,
		DisplayName: fmt.Sprintf("文件管理系统 %s", pkg.AppVersion),
		Description: "a file manager system",
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
	return boot(cfg)
}

func (m *GFSService) GetAny(s string) ([]byte, []string) {
	z.L().Sugar().Debugln("GFSService service GetAny")
	return input(), nil
}
