package util

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

const (
	DefaultInstallPath = "/usr/local/" + MarketName
)

// GetDiskUsage 获取 Unix 系统磁盘使用情况
func GetDiskUsage(path string) (total, used, free uint64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return
	}
	total = stat.Blocks * uint64(stat.Bsize)
	free = stat.Bfree * uint64(stat.Bsize)
	used = total - free
	return
}

func SetPlatformSpecificAttrs(cmd *exec.Cmd) {
	if runtime.GOOS == "linux" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true, // 创建新会话，脱离终端
			//Setpgid: true, // 创建新的进程组
			//Pgid:    0,    // 子进程成为进程组领导者
			// 或者使用 Setsid: true 创建新会话（类似 nohup）
		}
		// 重定向输入输出（避免挂起）
		//cmd.Stdout = nil
		//cmd.Stderr = nil
		//cmd.Stdin = nil
	}
}

func getOsName() (osName string) {
	if runtime.GOOS == "android" {
		return "Android"
	}
	var sysnamePath string
	sysnamePath = "/etc/redhat-release"
	_, err := os.Stat(sysnamePath)
	if err != nil && os.IsNotExist(err) {
		str := "PRETTY_NAME="
		f, err := os.Open("/etc/os-release")
		if err != nil && os.IsNotExist(err) {
			str = "DISTRIB_ID="
			f, err = os.Open("/etc/openwrt_release")
		}
		if err == nil {
			buf := bufio.NewReader(f)
			for {
				line, err := buf.ReadString('\n')
				if err == nil {
					line = strings.TrimSpace(line)
					pos := strings.Count(line, str)
					if pos > 0 {
						len1 := len([]rune(str)) + 1
						rs := []rune(line)
						osName = string(rs[len1 : (len(rs))-1])
						break
					}
				} else {
					break
				}
			}
		}
	} else {
		buff, err := ioutil.ReadFile(sysnamePath)
		if err == nil {
			osName = string(bytes.TrimSpace(buff))
		}
	}
	if osName == "" {
		osName = "Linux"
	}
	return
}

func SetRLimit() error {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return err
	}
	limit.Max = 65536
	limit.Cur = limit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return err
	}
	return nil
}

func SetFirewall(ProductName, fullPath string) {
}
