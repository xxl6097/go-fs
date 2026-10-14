package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/utils/util"
)

func RestartWindowsApplication() {
	if IsWindows() {
		// Windows特有重启逻辑
		exe, _ := os.Executable()
		cmd := exec.Command("cmd", "/C", "timeout 2 && "+exe)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			z.Error(err)
		}
		os.Exit(0)
	}
}

func BlockingFunction[T any](c context.Context, timeout time.Duration, callback func() T) (T, error) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	resultChan := make(chan T)
	go func() {
		result := callback()
		resultChan <- result
	}()
	var zero T // 声明 T 的零值
	select {
	case res := <-resultChan:
		return res, nil
	case <-ctx.Done():
		return zero, errors.New("timeout")
	}
}
func DynamicSelect[T any](t []T, fun func(context.Context, int, T) T) T {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan T, len(t)) // 缓冲大小等于协程数量
	var wg sync.WaitGroup
	for i, v := range t {
		wg.Add(1)
		go func(ct context.Context, index int, t T, c chan<- T) {
			defer wg.Done()
			c <- fun(ct, index, t)
		}(ctx, i, v, ch)
	}
	var x T
	for i := 0; i < len(t); i++ {
		_, value, ok := reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}})
		r := value.Interface().(T)
		if ok {
			cancel()
			wg.Wait()
			return r
		}
	}
	cancel()
	return x
}
func IsMacOs() bool {
	if strings.Compare(runtime.GOOS, "darwin") == 0 {
		return true
	}
	return false
}
func IsWindows() bool {
	if strings.Compare(runtime.GOOS, "windows") == 0 {
		return true
	}
	return false
}
func IsOpenWRT() bool {
	_, err := os.Stat("/etc/openwrt_release")
	if err == nil {
		return true
	}
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}
	src := strings.ToLower(string(data))
	dst := strings.ToLower("OpenWRT")
	return strings.Contains(src, dst)
}

func RunChildProcess(executable string, args ...string) error {
	// 在 systemd 环境下，若直接 fork 子进程去执行 restart，子进程会与当前服务
	// 处于同一个 service cgroup；一旦 restart 触发 systemctl stop，systemd 默认的
	// KillMode=control-group 会把整个 cgroup（含这个正在执行 restart 的子进程）一起
	// SIGTERM，导致「stop 成功、start 未完成」，服务停在停止态。
	// 用 systemd-run 把子进程放进独立的 transient scope（独立 cgroup）即可规避连坐。
	if useSystemdRun() {
		runArgs := []string{
			"--scope",
			"--collect", // 运行结束后自动清理该 transient unit
			fmt.Sprintf("--unit=%s-restart-%d", filepath.Base(executable), time.Now().UnixNano()),
			executable,
		}
		runArgs = append(runArgs, args...)
		cmd := exec.Command("systemd-run", runArgs...)
		util.SetPlatformSpecificAttrs(cmd)
		fmt.Printf("运行子进程(systemd-run) %s %v\n", executable, args)
		if err := cmd.Start(); err == nil {
			return nil
		} else {
			// systemd-run 启动失败不能回退到普通 fork：普通 fork 的子进程仍在这个
			// service cgroup 里，restart 触发 systemctl stop 时会被连带 SIGTERM，
			// restart 照样不生效。直接报错，让调用方知道重启没能被安全拉起。
			return fmt.Errorf("systemd-run 启动失败: %w", err)
		}
	}
	// systemd 环境但缺少 systemd-run：普通 fork 逃不出 service cgroup，同上会连坐。
	if isSystemd() {
		return errors.New("systemd 环境缺少 systemd-run，无法安全拉起重启子进程")
	}
	cmd := exec.Command(executable, args...)
	util.SetPlatformSpecificAttrs(cmd)
	fmt.Printf("运行子进程 %s %v\n", executable, args)
	return cmd.Start()
}

// isSystemd 判断是否运行在 systemd 管理的 Linux 系统上。
func isSystemd() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	// systemd 运行时目录存在即视为 systemd 系统
	if _, err := os.Stat("/run/systemd/system"); err != nil {
		return false
	}
	return true
}

// useSystemdRun 判断当前是否为可用 systemd-run 的 Linux systemd 环境。
func useSystemdRun() bool {
	if !isSystemd() {
		return false
	}
	if _, err := exec.LookPath("systemd-run"); err != nil {
		return false
	}
	return true
}

func RunCmdBySelf(name string, args ...string) error {
	//defer glog.Flush()
	//binpath, err := os.Executable()
	//if err != nil {
	//	return err
	//}
	z.Debug("RunCmdBySelf", name, args)
	err := RunChildProcess(name, args...)
	if err != nil {
		z.Errorf("RunChildProcess错误: %v\n", err)
		return fmt.Errorf("RunChildProcess错误: %v\n", err)
	}
	z.Println("子进程启动成功", name)
	return err
}

func RunCmdWithSudo(args ...string) ([]byte, error) {
	z.Debug("run", args)
	cmd := exec.Command("sudo", args...)
	output, err := cmd.CombinedOutput() // 捕获标准输出和错误
	if err != nil {
		return nil, err
	}
	fmt.Println(string(output)) // 输出：hello world
	return output, err

	//cmd := exec.Command("sudo", args...)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//// 执行命令
	//return cmd.Run()
}

func RunCmd(name string, args ...string) string {
	output, err := Cmd(name, args...)
	if err != nil {
		return err.Error()
	}
	return string(output)
}

func Cmd(name string, args ...string) ([]byte, error) {
	z.L().Sugar().Debug("run", name, args)
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput() // 捕获标准输出和错误
	if err != nil {
		return nil, err
	}
	return output, err
}

func RunWithSudo() error {
	if os.Geteuid() == 0 {
		return nil // 已经拥有 root 权限
	}

	// 获取当前可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 构建 sudo 命令
	cmd := exec.Command("sudo", exePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}

func RunWithSudo1() error {
	if os.Geteuid() == 0 {
		return nil // 已经拥有 root 权限
	}

	// 获取当前可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 构建 sudo 命令
	cmd := exec.Command("sudo", exePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}

// PrintByteArrayAsConstant 把字节数组以常量字节数组的形式打印出来
func PrintByteArrayAsConstant(bytes []byte) string {
	sb := strings.Builder{}
	sb.WriteString("[]byte{")
	for i, b := range bytes {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("0x%02X", b))
	}
	sb.WriteString("}")
	return sb.String()
}

func ExitAnyKey() {
	fmt.Print("按回车键退出程序...")
	endKey := make([]byte, 1)
	_, _ = os.Stdin.Read(endKey) // 等待用户输入任意内容后按回车
	os.Exit(0)
}

func ExitCountDown(count int) {
	//for i := count; i >= 0; i-- {
	//	fmt.Printf("\r%d秒后退出程序..", i)
	//	time.Sleep(1 * time.Second)
	//}
	fmt.Println("✅ 安装完成，退出！")
	os.Exit(0)
}

// GzipCompress 压缩字节数组
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	// 设置压缩级别（BestSpeed, BestCompression, DefaultCompression）
	gz, e := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if e != nil {
		return nil, e
	}
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil { // 必须关闭以写入所有数据
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipDecompress 解压
func GzipDecompress(compressed []byte) ([]byte, error) {
	r, e := gzip.NewReader(bytes.NewReader(compressed))
	if e != nil {
		return nil, e
	}
	defer r.Close()
	return io.ReadAll(r)
}

func GetSelfSize() uint64 {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行文件路径时出错: %v\n", err)
		return 0
	}
	// 获取文件信息
	fileInfo, err := os.Stat(exePath)
	if err != nil {
		fmt.Printf("获取文件信息时出错: %v\n", err)
		return 0
	}

	// 获取文件大小
	fileSize := fileInfo.Size()
	fmt.Printf("本程序自身大小为: %v\n", ByteCountIEC(uint64(fileSize)))
	return uint64(fileSize)
}
