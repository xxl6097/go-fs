package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xxl6097/gfs/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/gofs"
	"go.uber.org/zap"
)

// Upgrader 实现 gofs.Upgrader。
type Upgrader struct {
	gs igs.Service
}

func (u *Upgrader) Check(ctx context.Context) (*gofs.Release, error) {
	var releaseNote *gofs.Release
	github.Api().SetName("xxl6097", "go-fs")
	data, err := github.Api().CheckUpgrade(pkg.BinName)
	if err != nil {
		// 已是最新版时 err 为 "已是最新版本～"
		z.L().Sugar().Warnln(err)
		return releaseNote, nil
	}
	publishedAt := data["publishedAt"].(time.Time) // markdown 文本，可直接渲染
	releaseNote = &gofs.Release{
		Version:     data["version"].(string),
		Notes:       data["releaseNotes"].(string),
		PublishedAt: publishedAt.Format(time.DateTime),
		PageURL:     data["htmlUrl"].(string),
		FullURL:     data["fullUrl"].(string),
		PatchURL:    data["patchUrl"].(string),
	}
	return releaseNote, err
}
func (u *Upgrader) Upgrade(ctx context.Context, rel *gofs.Release, onProgress func(gofs.UpgradeProgress)) error {
	urls := github.Api().GetProxyUrls(rel.FullURL)
	fileUri := utils.DownloadFileWithCancelByUrls(urls, func(progress float64) {

	}) // 并发竞速，返回 "" 表示失败
	if fileUri != "" {
		z.L().Debug("升级文件", zap.String("fileUri", fileUri))
		err := u.gs.Upgrade(ctx, fileUri)
		return err
	}

	z.L().Debug("文件地址空fileUri", zap.String("fileUri", fileUri))
	return fmt.Errorf("文件地址空fileUri=%s", fileUri)
}

func NewUpgrader(gs igs.Service) *Upgrader {
	return &Upgrader{gs: gs}
}
