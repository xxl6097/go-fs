package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/xxl6097/glog/pkg/z"
)

//func getSegmentValue(seg []string, idx int) int {
//	if idx >= len(seg) {
//		return 0 // 自动补零处理长度不一致情况
//	}
//	num, _ := strconv.Atoi(seg[idx])
//	return num
//}
//
//func SplitVersion(v string) []string {
//	// 去除前缀标识（如 "v1.2.3" → "1.2.3"）
//	v = strings.TrimLeft(v, "v")
//	return strings.Split(v, ".")
//}

// CompareVersionsExt 0:相等；1：v1>v2;-1:v1<v2
func CompareVersionsExt(v1, v2 string) int {
	seg1 := SplitVersion(v1)
	seg2 := SplitVersion(v2)
	maxLen := int(math.Max(float64(len(seg1)), float64(len(seg2))))

	for i := 0; i < maxLen; i++ {
		num1 := getSegmentValue(seg1, i)
		num2 := getSegmentValue(seg2, i)

		if num1 > num2 {
			return 1 // v1 > v2
		} else if num1 < num2 {
			return -1 // v1 < v2
		}
	}
	return 0 // 相等
}

// CompareVersions 0:相等；大于0有新版本，小于零无新版本
func CompareVersions1(new, old string) int {
	seg1 := SplitVersion(new)
	seg2 := SplitVersion(old)
	maxLen := int(math.Max(float64(len(seg1)), float64(len(seg2))))

	var data1, data2 int
	for i := 0; i < maxLen; i++ {
		num1 := getSegmentValue(seg1, maxLen-i-1)
		num2 := getSegmentValue(seg2, maxLen-i-1)
		fang := int(math.Pow(100, float64(i)))
		data1 += num1 * fang
		data2 += num2 * fang
		z.Debug(num1, num2, fang, seg1, seg2)
		z.Debug(data1, data2)
	}
	return data1 - data2
}

//func CompareVersions(new, old string) int {
//	a := strings.ReplaceAll(new, "v", "")
//	a = strings.ReplaceAll(a, ".", "")
//	data1, err := strconv.ParseInt(a, 10, 64)
//	if err != nil {
//		return CompareVersions1(new, old)
//	}
//	b := strings.ReplaceAll(old, "v", "")
//	b = strings.ReplaceAll(b, ".", "")
//	data2, err := strconv.ParseInt(b, 10, 64)
//	if err != nil {
//		return CompareVersions1(new, old)
//	}
//	return int(data1 - data2)
//}

func GetVersionByFileName(filename string) string {
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	//fmt.Println(re.FindStringSubmatch(filename))
	return re.FindString(filename)
}

func ReplaceNewVersionBinName(filename, v string) string {
	re := regexp.MustCompile(`_v\d+\.\d+\.\d+_`)
	newName := re.ReplaceAllString(filename, fmt.Sprintf("_%s_", v)) // 替换为单个下划线
	fmt.Println(newName)
	return newName
}

func ParseMarkdownCodeToStringArray(body string) []string {
	codeBlocks := ExtractCodeBlocks(body)
	if len(codeBlocks) == 0 {
		return []string{}
	}
	var r []string
	err := json.Unmarshal([]byte(codeBlocks[0]), &r)
	if err != nil {
		return []string{}
	}
	return r
}

func SplitVersion(v string) []string {
	// 去除前缀标识（如 "v1.2.3" → "1.2.3"）
	v = strings.TrimLeft(v, "v")
	return strings.Split(v, ".")
}

// CompareVersions 0:相等；1：v1>v2;-1:v1<v2
func CompareVersions(v1, v2 string) int {
	seg1 := SplitVersion(v1)
	seg2 := SplitVersion(v2)
	maxLen := int(math.Max(float64(len(seg1)), float64(len(seg2))))

	for i := 0; i < maxLen; i++ {
		num1 := getSegmentValue(seg1, i)
		num2 := getSegmentValue(seg2, i)

		if num1 > num2 {
			return 1 // v1 > v2
		} else if num1 < num2 {
			return -1 // v1 < v2
		}
	}
	return 0 // 相等
}

func getSegmentValue(seg []string, idx int) int {
	if idx >= len(seg) {
		return 0 // 自动补零处理长度不一致情况
	}
	num, _ := strconv.Atoi(seg[idx])
	return num
}
