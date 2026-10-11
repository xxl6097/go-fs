package main

import (
	"fmt"
	"log"

	"github.com/xxl6097/go-service/pkg/github"
)

func main() {
	github.Api().SetName("xxl6097", "go-fs")

	data, err := github.Api().CheckUpgrade("gfs_v0.0.8_linux_amd64")
	if err != nil {
		// 已是最新版时 err 为 "已是最新版本～"
		log.Println(err)
		return
	}
	fmt.Println(data)
	// 有新版：
	fullUrl := data["fullUrl"].(string)    // 全量包下载地址，可能为空
	patchUrl := data["patchUrl"].(string)  // 差量包下载地址，可能为空
	notes := data["releaseNotes"].(string) // markdown 文本，可直接渲染
	fmt.Println(notes)
	fmt.Println(patchUrl)
	fmt.Println(fullUrl)

	uts := github.Api().GetProxyUrls(fullUrl)
	fmt.Println(uts)

}
