package service

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/xxl6097/go-service/pkg/utils"
)

// pathCompleter 实现 readline.AutoCompleter，用于文件/目录路径的 Tab 补全。
type pathCompleter struct{}

// Do 根据光标前已输入的内容补全路径。
//
// 约定（与 readline.AutoCompleter 一致）：
//   - candidates 里每一项是「替换光标前 length 个字符」的后缀文本；
//   - length 是待补全名字部分的字符（rune）数。
//
// 例：输入 /home/uuxia/Deskt 时 dir=/home/uuxia/、partial=Deskt、length=5，
// 候选为匹配项去掉 Deskt 前缀后的剩余部分（目录再补一个 /）。
func (pathCompleter) Do(line []rune, pos int) (candidates [][]rune, length int) {
	input := string(line[:pos])
	dir, partial := splitPath(input)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0
	}
	for _, e := range entries {
		name := e.Name()
		// 空输入时不提示隐藏文件（与 shell 习惯一致，避免 .git/.DS_Store 刷屏）
		if partial == "" && strings.HasPrefix(name, ".") {
			continue
		}
		if !strings.HasPrefix(name, partial) {
			continue
		}
		suffix := name[len(partial):]
		if e.IsDir() {
			suffix += string(os.PathSeparator)
		}
		candidates = append(candidates, []rune(suffix))
	}
	return candidates, len([]rune(partial))
}

// splitPath 把输入拆成「目录部分」和「待补全的名字部分」。
// 例：/home/uuxia/Deskt -> (/home/uuxia/, Deskt)；Desktop -> (., Desktop)。
func splitPath(input string) (dir, partial string) {
	input = expandHome(input)
	i := strings.LastIndexAny(input, `/\`)
	if i < 0 {
		// 没有分隔符：补全当前目录下的名字
		return ".", input
	}
	dir, partial = input[:i], input[i+1:]
	if dir == "" {
		// 以 / 开头：补全根目录
		dir = string(os.PathSeparator)
	}
	return dir, partial
}

// expandHome 把路径开头的 ~ 展开成用户主目录。
func expandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return home + path[1:]
		}
	}
	return path
}

// InputPath 提示用户输入一个路径，终端下支持 Tab 补全；输入为空会一直重试。
// 非终端环境（stdin 不是 tty）退回普通行输入（无 Tab 补全）。
func InputPath(title string) (string, error) {
	if !readline.DefaultIsTerminal() {
		return utils.InputString(title), nil
	}
	for {
		line, err := readPathLine(title)
		if err != nil {
			return "", err
		}
		if line != "" {
			return line, nil
		}
		fmt.Println("【路径】不允许输入空")
	}
}

// readPathLine 读取一行路径输入，返回去除空白并展开 ~ 的结果。
func readPathLine(title string) (string, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          title,
		AutoComplete:    pathCompleter{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	line, err := l.Readline()
	if err != nil {
		if err == io.EOF || err == readline.ErrInterrupt {
			// Ctrl-D / Ctrl-C：中断输入
			return "", err
		}
		return "", err
	}
	return strings.TrimSpace(expandHome(line)), nil
}
