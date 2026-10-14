package zutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var serviceName = "aalog"

func Register(name string) {
	serviceName = name
}

func NameByPath(appPath string) string {
	appName := filepath.Base(appPath)
	// 处理通用扩展名
	if ext := filepath.Ext(appName); ext != "" {
		appName = strings.TrimSuffix(appName, ext)
	}
	if strings.Contains(appName, "_") {
		arr := strings.Split(appName, "_")
		if arr != nil && len(arr) > 0 {
			if arr[0] != "" {
				appName = arr[0]
			}
		}
	}
	if strings.Contains(appName, "-") {
		arr := strings.Split(appName, "-")
		if arr != nil && len(arr) > 0 {
			if arr[0] != "" {
				appName = arr[0]
			}
		}
	}
	if strings.Contains(appName, ".") {
		arr := strings.Split(appName, ".")
		if arr != nil && len(arr) > 0 {
			appName = arr[0]
			if arr[0] != "" {
				appName = arr[0]
			}
		}
	}
	return appName
}
func AppName() string {
	var appPath string
	if exePath, err := os.Executable(); err == nil {
		appPath = exePath
	} else {
		appPath = os.Args[0]
	}
	return NameByPath(appPath)
}
func TempDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("ProgramData"))
	default:
		return filepath.Join(os.TempDir())
	}
}
func AppHome(args ...string) string {
	var dirs []string
	dirs = append(dirs, TempDir(), serviceName)
	dirs = append(dirs, AppName())
	if len(args) > 0 {
		dirs = append(dirs, args...)
	}
	dir := filepath.Join(dirs...)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	return dir
}

func TextToFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", content)) // 直接写入字符串
	// 或写入字节切片：file.Write([]byte("data"))
}

func Now() time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	now := time.Now()
	beijingTime := now.In(loc)
	return beijingTime
}
