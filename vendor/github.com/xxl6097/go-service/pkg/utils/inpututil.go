package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func tips(title string) {
	str := strings.ReplaceAll(title, "请输入", "")
	str = strings.ReplaceAll(str, "please input", "")
	str = strings.ReplaceAll(str, "：", "")
	str = strings.ReplaceAll(str, ":", "")
	str = fmt.Sprintf("【%s】不允许输入空", str)
	fmt.Println(str)
}
func InputStringEmpty1(title string) string {
	reader := bufio.NewReader(os.Stdin)
	//glog.Print(title)
	fmt.Print(title)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return InputString(title)
	}
	//return strings.TrimSpace(input)
	return input
}
func InputStringEmpty(title, defaultString string) string {
	reader := bufio.NewReader(os.Stdin)
	//glog.Print(title)
	fmt.Print(title)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return InputString(title)
	}
	if input == "" {
		return defaultString
	}
	//return strings.TrimSpace(input)
	return input
}

func InputString(title string) string {
	reader := bufio.NewReader(os.Stdin)
	//glog.Print(title)
	fmt.Print(title)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return InputString(title)
	}
	//return strings.TrimSpace(input)
	if len(input) == 0 {
		tips(title)
		return InputString(title)
	}
	return input
}
func InputIntDefault(title string, def int) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(title)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return InputInt(title)
	}
	var value int
	if len(input) == 0 {
		return def
	} else {
		value, err = strconv.Atoi(input)
		if err != nil {
			tips(title)
			return InputInt(title)
		}
	}
	return value
}

func InputInt(title string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(title)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return InputInt(title)
	}
	if len(input) == 0 {
		tips(title)
		return InputInt(title)
	}
	num, err := strconv.Atoi(input)
	if err != nil {
		return InputInt(title)
	}
	return num
}

func GetInt() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if err != nil {
		return GetInt()
	}
	if len(input) == 0 {
		fmt.Println("不允许输入空")
		return GetInt()
	}
	num, err := strconv.Atoi(input)
	if err != nil {
		return GetInt()
	}
	return num
}
