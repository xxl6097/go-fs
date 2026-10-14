package model

import (
	"github.com/xxl6097/go-service/pkg/utils"
	"path"
	"path/filepath"
	"strings"
)

type Node struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	FilePath string `json:"filePath"`
}
type Option struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Children []Node `json:"children"`
}

func SplitLastTwoByUnderscore(s string) []string {
	// 过滤空元素
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_'
	})
	if len(parts) < 2 {
		return []string{}
	}
	return parts[len(parts)-2:]
}

func ToTree(dir string, entries []string) []Option {
	maps := make(map[string][]Node)
	for _, uri := range entries {
		var name string
		var fileUri string
		if utils.IsURL(uri) {
			name = path.Base(uri) // 输出：report.pdf
			fileUri = uri
		} else {
			name = uri
			fileUri = filepath.Join(dir, name)
		}
		result := SplitLastTwoByUnderscore(name)
		//fmt.Printf("%-30s => %v\n", name, result)
		if len(result) == 2 {
			nodeArray := maps[result[0]]
			if nodeArray == nil {
				nodeArray = make([]Node, 0)
			}
			nameOnly := utils.CleanExt(result[1])

			nodeArray = append(nodeArray, Node{
				Label:    nameOnly,
				Value:    nameOnly,
				FilePath: fileUri,
			})
			maps[result[0]] = nodeArray
		}
	}

	var options []Option
	for k, v := range maps {
		options = append(options, Option{
			Label:    utils.ToUpperFirst(k),
			Value:    k,
			Children: v,
		})
	}
	return options
}
