package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/go-service/pkg/utils/util"
)

func init() {
	InitLog()
}

type CoreService struct {
	isrv    igs.IService
	srv     service.Service
	config  *service.Config
	workDir string
}

func InitLog() {
	//glog.Register(util.MarketName)
	//bindir, err := os.Executable()
	//var isSrvApp bool
	//if err == nil {
	//	isSrvApp = strings.HasPrefix(strings.ToLower(bindir), strings.ToLower(util.DefaultInstallPath))
	//	if !isSrvApp {
	//		//fmt.Println("非服务程序。。。")
	//		//glog.SetLogFile(filepath.Dir(bindir), fmt.Sprintf("install-%s.log", filepath.Base(bindir)))
	//		z.LoadLogger(func(conf *z.LogConfig) {
	//			conf.Path = filepath.Join(filepath.Dir(bindir), fmt.Sprintf("install-%s.log", filepath.Base(bindir)))
	//		})
	//		return
	//	}
	//}
	//glog.LogDefaultLogSettingEveryType("app.log", everyType)
	z.LoadLogger(func(cfg *z.LogConfig) {
		cfg.LogDir = zutil.AppHome("log")
		cfg.MaxSize = 5
	})
}

func IsLogDirExist() bool {
	logDir := zutil.AppHome("log")
	info, err := os.Stat(logDir)
	fmt.Println("日志路径", logDir, info, err)
	// 1. 不存在：返回 false, nil
	if os.IsNotExist(err) {
		return false
	}
	// 2. 其他错误（权限不足等）
	if err != nil {
		return false
	}
	// 3. 存在，但要判断是不是文件夹
	return info.IsDir()
}

func (this *CoreService) Run() error {
	//defer glog.GlobalRecover()
	if len(os.Args) > 1 && os.Args[1] == utils.SYSTEM_CPU_INFO {
		fmt.Println("hello world")
		fmt.Println(this.isrv.OnVersion())
		return nil
	}
	this.config = this.isrv.OnConfig()
	if this.config == nil {
		return errors.New("请设置服务配置信息～")
	}
	if this.config.Name == "" {
		return errors.New("【Name】应用服务名不能为空字符串")
	}
	if this.config.DisplayName == "" {
		return errors.New("【DisplayName】应用服务显示名不能为空字符串")
	}
	if this.config.Description == "" {
		return errors.New("【Description】应用服务显示名不能为空字符串")
	}
	this.config.Name = strings.ToLower(this.config.Name)
	this.workDir = filepath.Join(util.DefaultInstallPath, this.config.Name)
	if utils.IsWindows() {
		this.config.Name = fmt.Sprintf("%s.exe", this.config.Name)
	}
	this.config.Executable = filepath.Join(this.workDir, this.config.Name)

	this.deleteOld()
	binDir := filepath.Dir(os.Args[0])
	_ = os.Chdir(binDir)
	e := this.createService()
	if e != nil {
		return e
	}
	//glog.Debug("运行参数", os.Args, os.Getpid())
	return this.menu()
}

func (this *CoreService) reqeustWindowsUser() {
	if utils.IsWindows() {
		if this.config.UserName == "" {
			username := utils.InputStringEmpty("请输入windows登录用户名：", "")
			if username != "" {
				this.config.UserName = username
			}
			password := utils.InputStringEmpty("请输入windows登录用户密码：", "")
			if password != "" {
				if this.config.Option == nil {
					this.config.Option = make(map[string]interface{})
				}
				this.config.Option["Password"] = password
			}

		}
	}
}

// 在安装目录且服务处于运行状态
func (this *CoreService) isServiceApp() bool {
	binPath, err := os.Executable()
	if err != nil {
		return false
	}
	if strings.Compare(strings.ToLower(this.config.Executable), strings.ToLower(binPath)) == 0 {
		s, e := this.statusService()
		if e != nil {
			return false
		}
		if s == service.StatusRunning {
			return true
		}
	}
	return false
}

func (this *CoreService) deleteOld() {
	tempFilePath := filepath.Join(this.workDir, fmt.Sprintf(".%s.old", this.config.Name))
	if utils.FileExists(tempFilePath) {
		_ = utils.DeleteAllDirector(tempFilePath)
	}
}

func NewCore(srv igs.IService) igs.DefaultService {
	return &CoreService{isrv: srv}
}
