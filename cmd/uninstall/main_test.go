package main

import (
	"strings"
	"testing"

	"github.com/xxl6097/gofs"
)

// TestUninstallSatisfiesInterface 是「这个值能直接塞进 WithUninstaller」
// 这条契约的编译期版本。
func TestUninstallSatisfiesInterface(t *testing.T) {
	var _ gofs.Uninstaller = gofs.UninstallFunc(options{}.uninstall)
}

// TestNewOptionsExitMode 覆盖 -exit 的取值。
//
// 选错的后果是「点了之后页面既不跳也不报错」或「响应永远发不出去」，
// 都要手点一次才看得出来，所以两个分支连同未知取值一起钉住。
func TestNewOptionsExitMode(t *testing.T) {
	now, err := newOptions("now", "", false, false)
	if err != nil || !now.exitNow {
		t.Errorf("now 应当置上 exitNow：%+v, %v", now, err)
	}
	later, err := newOptions("later", "", false, false)
	if err != nil || later.exitNow {
		t.Errorf("later 不该置上 exitNow：%+v, %v", later, err)
	}
	if _, err := newOptions("", "", false, false); err == nil {
		t.Error("空 -exit 应当报错")
	}
	if _, err := newOptions("immediately", "", false, false); err == nil {
		t.Error("未知的 -exit 应当报错")
	}
}

// TestStepsDefaultTouchesNothing 是安全底线：**默认一个文件都不删**。
//
// 示例跑起来不该把使用者刚编译出来的二进制删掉。真有东西被删，
// 必须是显式给了开关 —— 这条测试守的就是这个默认值。
func TestStepsDefaultTouchesNothing(t *testing.T) {
	o, err := newOptions("later", "", false, false)
	if err != nil {
		t.Fatal(err)
	}
	steps := o.steps()
	if len(steps) != 1 {
		t.Fatalf("默认应当只剩「停止服务」一步，得到 %d 步：%v", len(steps), descs(steps))
	}
	if !strings.Contains(steps[0].desc, "停止服务") {
		t.Errorf("默认那一步应当是停止服务，得到 %q", steps[0].desc)
	}
	if !strings.Contains(steps[0].desc, "不会再被拉起来") {
		t.Errorf("提示里要说清「卸载不是重启」：%q", steps[0].desc)
	}
}

// TestStepsIncludeOnlyWhatWasAsked 确认清单跟着开关走，且顺序正确。
func TestStepsIncludeOnlyWhatWasAsked(t *testing.T) {
	o, err := newOptions("later", "/etc/systemd/system/gofs.service", true, true)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(descs(o.steps()), "\n")
	for _, want := range []string{"停止服务", "gofs.service", "配置目录", "二进制"} {
		if !strings.Contains(got, want) {
			t.Errorf("清单里缺少 %q：\n%s", want, got)
		}
	}

	// 顺序：先停服务、再删注册项、然后清配置目录、最后删二进制。
	//
	// 反过来（先删文件后停服务）会让服务处于「半死」状态：端口还占着、
	// 配置已经没了；先删二进制再删注册项则会留下一堆「启动失败」的告警。
	order := []string{"停止服务", "gofs.service", "配置目录", "二进制"}
	idx := 0
	for _, d := range descs(o.steps()) {
		if idx < len(order) && strings.Contains(d, order[idx]) {
			idx++
		}
	}
	if idx != len(order) {
		t.Errorf("步骤顺序不对，走到了第 %d 项：\n%s", idx, got)
	}

	// 只给其中一项时，别的不能顺带出现。
	one, err := newOptions("later", "", true, false)
	if err != nil {
		t.Fatal(err)
	}
	gotOne := strings.Join(descs(one.steps()), "\n")
	if strings.Contains(gotOne, "二进制") || strings.Contains(gotOne, "daemon-reload") {
		t.Errorf("没要求的步骤不该出现：\n%s", gotOne)
	}
}

// TestExitNote 只是确认两种模式都会给人一句可读的说明 —— 启动横幅要打印它。
func TestExitNote(t *testing.T) {
	now, _ := newOptions("now", "", false, false)
	later, _ := newOptions("later", "", false, false)
	if !strings.Contains(now.exitNote(), "断") {
		t.Errorf("now 的说明没提到连接会断：%q", now.exitNote())
	}
	if !strings.Contains(later.exitNote(), "200") {
		t.Errorf("later 的说明没提到响应会正常返回：%q", later.exitNote())
	}
}

func descs(steps []step) []string {
	out := make([]string, 0, len(steps))
	for _, s := range steps {
		out = append(out, s.desc)
	}
	return out
}
