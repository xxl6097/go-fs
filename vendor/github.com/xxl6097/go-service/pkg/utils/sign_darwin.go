//go:build darwin

package utils

import "os/exec"

// AdhocSign 对二进制做 ad-hoc 重签名。
//
// Go 链接器为 darwin/arm64 生成的可执行文件自带 ad-hoc 签名（linker-signed），
// 该签名覆盖整个文件。安装时 GenerateBin 会原地替换特征码写入配置，
// 文件被改写后签名即失效，Apple Silicon 内核会以 OS_REASON_CODESIGNING
// 直接 SIGKILL 掉进程（launchctl 里表现为退出状态 -9、反复重启）。
// 所以改写之后必须重新签名；ad-hoc 不需要证书和开发者账号。
func AdhocSign(filePath string) error {
	return exec.Command("codesign", "--force", "--sign", "-", filePath).Run()
}
