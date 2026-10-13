// Command uninstall-demo 演示怎么把 gofs 的「卸载」接起来。
//
// gofs 自己不实现「怎么卸载自己」—— 它只定义接口 gofs.Uninstaller，
// 实现由部署方给。为什么：卸载是比升级（换掉自己的二进制）更不可逆的操作，
// 该删什么、删到什么程度完全取决于部署方式（systemd / launchd / Windows
// 服务 / 容器编排），库替调用方拍板只会错。
//
// 用法：
//
//	go build -o /tmp/gofs-uninstall-demo ./examples/uninstall
//	/tmp/gofs-uninstall-demo -root ./data -port 8080
//
//	# 打开 http://127.0.0.1:8080 ，用 admin / admin123 登录，
//	# 进「设置 → 升级、重启与卸载」，点「卸载」。
//
// ⚠️ 默认**什么都不删**：本示例只演示流程，以及「进程退出 → 页面判定
// 卸载完成」这条链路。要真的动手，显式给出目标（见下面三个开关）：
//
//	/tmp/gofs-uninstall-demo -unit /etc/systemd/system/gofs.service \
//	    -purge-config -remove-binary
//
// ⚠️ 不要用 go run 跑它：os.Executable() 拿到的是编译出来的临时文件，
// -remove-binary 删的会是那个临时文件，没什么意义。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/xxl6097/gofs"
)

func main() {
	root := flag.String("root", ".", "要服务的目录")
	port := flag.Int("port", 8080, "监听端口")
	exit := flag.String("exit", "later",
		"何时退出进程：later＝响应发出去之后再退（页面拿到 200）；"+
			"now＝回调里立刻 os.Exit、不写响应（页面看到连接被掐断）")
	unit := flag.String("unit", "",
		"要删除的服务注册文件（如 /etc/systemd/system/gofs.service）；留空＝跳过这一步")
	purgeConfig := flag.Bool("purge-config", false,
		"是否删除 gofs 的配置目录（keys.json / users.json / settings.json）")
	removeBinary := flag.Bool("remove-binary", false,
		"是否删除自己的二进制")
	flag.Parse()

	opts, err := newOptions(*exit, *unit, *purgeConfig, *removeBinary)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, err := gofs.New(
		gofs.WithRoot(*root),
		gofs.WithBind("127.0.0.1"),
		gofs.WithPort(*port),
		gofs.WithAllowAll(),
		gofs.WithAuth("admin:admin123@/:rw"),
		// 密钥 / 用户表 / 运行期设置都放内存，别碰开发者本机的真实配置。
		gofs.WithKeyFile(""),
		gofs.WithUserFile(""),
		gofs.WithSettingsFile(""),
		// ★ 关键的一行：注入之后页面上才会出现「卸载」。
		// 不注入的话接口返回 404，按钮也不显示 —— 免得给一个点了必然失败的按钮。
		gofs.WithUninstaller(gofs.UninstallFunc(opts.uninstall)),
	)
	if err != nil {
		log.Fatalf("构造服务失败: %v", err)
	}
	defer srv.Close()

	if err := srv.Start(); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
	srv.Banner(os.Stdout)
	fmt.Printf("\n打开 %s ，用 admin / admin123 登录，\n", srv.URL())
	fmt.Printf("在「设置 → 升级、重启与卸载」里点「卸载」。\n")
	fmt.Printf("本次会做的事（%s）：\n", opts.exitNote())
	for _, s := range opts.steps() {
		fmt.Printf("  · %s\n", s.desc)
	}
	fmt.Println()

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
	log.Printf("服务已停止")
}

// options 是这次卸载要做什么。
//
// 拆成结构体而不是直接把 flag 传进回调，是为了让「会删哪些东西」能被列出来、
// 也能被测试：卸载是本服务最不可逆的操作，动手之前先把清单打出来，
// 免得点下去才发现删多了。命令行启动时会把清单打印一遍。
type options struct {
	// exitNow 为 true 时在回调里直接退出进程，false 则等响应发出去之后再退。
	exitNow bool
	// unit / purgeConfig / removeBinary 决定要不要真的删东西，
	// 默认全关 —— 示例不该把使用者刚编译出来的东西删掉。
	unit         string
	purgeConfig  bool
	removeBinary bool
}

func newOptions(exit, unit string, purgeConfig, removeBinary bool) (options, error) {
	o := options{unit: unit, purgeConfig: purgeConfig, removeBinary: removeBinary}
	switch exit {
	case "now":
		o.exitNow = true
	case "later":
		o.exitNow = false
	default:
		return options{}, fmt.Errorf("未知的 -exit %q，可选：now / later", exit)
	}
	return o, nil
}

func (o options) exitNote() string {
	if o.exitNow {
		return "-exit now：回调里立刻退出、不写响应 —— 页面会看到连接被掐断，\n" +
			"     这正是前端「把断连当成功」要处理的那一支"
	}
	return "-exit later：响应先发出去（页面拿到 200），进程随后退出"
}

func (o options) configDir() string { return gofs.DefaultKeyDir() }

// selfPath 返回自己的可执行文件路径。
func selfPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("定位当前可执行文件失败: %w", err)
	}
	// 解析软链只是为了报错信息好看（/usr/local/bin/gofs 常是软链）。
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}
	return exe, nil
}

// steps 按**正确的顺序**列出这次卸载要做的事。
//
// 顺序有讲究，别图省事调换：
//
//  1. 先停服务 —— 反过来的话，文件删了服务还在跑，处于「半死」状态：
//     端口还占着、日志还在写，而你已经没有它的配置了；
//  2. 再删注册项（systemd unit / launchd plist）—— 否则它开机还会被拉起来，
//     而那时二进制已经没了，留下一堆「启动失败」的告警；
//  3. 然后清配置目录；
//  4. 最后删二进制（非 Windows：删正在运行的可执行文件是可行的，
//     inode 会保留到进程退出）。
//
// 返回的每项都带 desc，调用方先把清单打出来再动手。
func (o options) steps() []step {
	out := []step{{desc: "停止服务（卸载不是重启，不会再被拉起来）"}}

	if o.unit != "" {
		unit := o.unit
		out = append(out, step{
			desc: "删除服务注册文件 " + unit + " 并 daemon-reload",
			run: func() error {
				if err := os.Remove(unit); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("删除 %s 失败: %w", unit, err)
				}
				// 删完 unit 要 reload，否则 systemd 内存里那份定义还在。
				if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
					// 不是 systemd 部署（比如 macOS）时报错很正常，不必拦下卸载。
					log.Printf("systemctl daemon-reload 未执行成（非 systemd 部署属正常）：%v", err)
				}
				return nil
			},
		})
	}

	if o.purgeConfig {
		dir := o.configDir()
		out = append(out, step{
			desc: "删除配置目录 " + dir,
			run: func() error {
				if err := os.RemoveAll(dir); err != nil {
					return fmt.Errorf("删除配置目录 %s 失败: %w", dir, err)
				}
				return nil
			},
		})
	}

	if o.removeBinary {
		out = append(out, step{
			desc: "删除自己的二进制",
			run: func() error {
				exe, err := selfPath()
				if err != nil {
					return err
				}
				if runtime.GOOS == "windows" {
					// Windows 不允许删除正在运行的可执行文件，
					// 真实部署这里通常交给一个「重启后删除」的批处理。
					log.Printf("Windows 上无法删除正在运行的二进制，跳过 %s", exe)
					return nil
				}
				return os.Remove(exe)
			},
		})
	}

	return out
}

// step 是卸载中的一步：先打日志（desc），再执行（run 为 nil 表示只报告不动作）。
type step struct {
	desc string
	run  func() error
}

// uninstall 交给 gofs 调用。它拿到的 ctx 已经脱开请求上下文并带 5 分钟超时
// （浏览器关掉页面不会中断卸载）。
//
// 返回 nil 表示「已经安排好了」—— 与 Restart 同一个契约：这个 HTTP 响应
// 是尽力而为。下面 -exit now 那一支就直接退出了，响应发不出去，页面只看到
// 连接被掐断，随后靠轮询 /fs/health **没**回来确认卸载成功。
func (o options) uninstall(ctx context.Context) error {
	log.Printf("收到卸载请求，开始执行（ctx 已脱开请求，超时 %s）", uninstallTimeoutNote)

	for _, s := range o.steps() {
		log.Printf("· %s", s.desc)
		if s.run == nil {
			continue
		}
		if err := s.run(); err != nil {
			// 返回错误 → 服务端回 500，页面把这条文案原样显示出来。
			// 比「假装成功」好得多：那时服务还在跑，管理员却以为已经卸干净了。
			return err
		}
	}

	log.Printf("卸载步骤执行完毕，准备停止服务")
	return o.stop()
}

// stop 停掉自己。
func (o options) stop() error {
	if o.exitNow {
		// 立刻退出，**不写响应**：客户端看到连接被掐断 —— 这正是前端
		// 「把断连当成卸载的正常表现」那一支要处理的情况。
		//
		// ⚠️ 注意这里用 os.Exit 而不是给自己发 SIGTERM：gofs 的 Run 装了
		// SIGTERM 处理器，会走**优雅关闭**，把在途请求（包括当前这个）
		// 等完再退 —— 那样响应照样发得出去，也就演示不出「断连」这条路径。
		//
		// 代价是跳过优雅关闭：在途的上传/下载会被掐断，日志与密钥也来不及
		// 刷盘。卸载场景下通常可以接受（配置目录本来就要删掉），
		// 但在乎在途传输的部署应当用 -exit later。
		os.Exit(0)
	}
	// 先让响应发出去，再退：handler 写完响应体之后这里返回，半秒后
	// 再给自己发 SIGTERM，走 gofs 的优雅关闭。
	go func() {
		time.Sleep(500 * time.Millisecond)
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()
	return nil
}

// uninstallTimeoutNote 与 internal/server 的 uninstallTimeout 保持一致。
// 示例不 import internal，所以这里写个常量说明。
const uninstallTimeoutNote = "5m"
