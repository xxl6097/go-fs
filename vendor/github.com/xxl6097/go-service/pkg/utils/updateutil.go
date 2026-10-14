package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-update"
	"go.uber.org/zap"
)

const SYSTEM_CPU_INFO = "system_cpu_info"

func PerformUpdate(newFilePath, targetPath string, patcher bool) error {
	file, err := os.Open(newFilePath)
	if err != nil {
		z.Error("升级文件打开失败", err)
		return fmt.Errorf("升级文件打开失败【%s】: %v", newFilePath, err)
	}
	defer func() {
		_ = file.Close()
	}()
	// Windows需要管理员权限
	opts := update.Options{
		TargetPath: targetPath, // 当前可执行文件路径
		Middler:    IsMissMatchOsApp,
	}
	if patcher {
		opts.Patcher = update.NewBSDiffPatcher()
	}

	//opts.CheckPermissions()
	//opts := update.Options{
	//	TargetPath: os.Args[0], // 当前可执行文件路径
	//}
	// 使用 bufio.NewReader 创建带缓冲的读取器
	z.L().Sugar().Debug("准备升级", patcher, newFilePath, targetPath)
	if err = update.Apply(bufio.NewReader(file), opts); err != nil {
		z.Error("升级失败，回滚", err)
		if e := update.RollbackError(err); e != nil {
			z.Error("更新失败且无法回滚", e)
			return fmt.Errorf("更新失败且无法回滚: %w", e)
		}
		return fmt.Errorf("apply失败: %v", err)
	}
	z.L().Sugar().Debug("升级执行完成", patcher, newFilePath, targetPath)
	return nil
}

func IsMissMatchOsApp(binPath string) error {
	if !FileExists(binPath) {
		z.L().Error("文件不存在", zap.String("binPath", binPath))
		return fmt.Errorf("文件不存在")
	}
	err := os.Chmod(binPath, 0755)
	if err != nil {
		return fmt.Errorf("赋予权限错误 %v", err)
	}
	o, e := Cmd(binPath, SYSTEM_CPU_INFO)
	if e != nil {
		//Copy(binPath, fmt.Sprintf("/tmp/aaaa_%d", time.Now().UnixMilli()))
		//_ = CopyToTemp(binPath)
		return fmt.Errorf("cmd运行错误 %s %v", binPath, e)
	}

	//z.Debug(binPath)
	z.L().Sugar().Debug("运行结果:", string(o))
	return nil
}

func ExtractCodeBlocks(markdown string) []string {
	var codeBlocks []string
	inCodeBlock := false
	var currentCodeBlock strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(markdown))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				codeBlocks = append(codeBlocks, currentCodeBlock.String())
				currentCodeBlock.Reset()
			}
			inCodeBlock = !inCodeBlock
		} else if inCodeBlock {
			currentCodeBlock.WriteString(line)
			currentCodeBlock.WriteRune('\n')
		}
	}

	return codeBlocks
}
