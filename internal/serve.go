package internal

import (
	"context"
	"fmt"
	"io"

	"github.com/xxl6097/gofs"
)

func RunCfg(ctx context.Context, cfg *gofs.Config, uploader gofs.Upgrader, restarter gofs.Restarter, stdout, stderr io.Writer) error {
	srv, err := gofs.New(gofs.WithConfig(cfg),
		gofs.WithUpgrader(uploader),
		gofs.WithRestarter(restarter))
	if err != nil {
		return err
	}
	defer srv.Close()

	if err := srv.Start(); err != nil {
		return err
	}
	srv.Banner(stdout)

	select {
	case <-ctx.Done():
		fmt.Fprintln(stderr, "\ngofs: 正在关闭…")
		shutdownCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx), gofs.DefaultShutdownTimeout)
		defer cancel()
		// 关闭超时不算运行失败：连接已被强制断开，提示一句就够了。
		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "gofs: %v\n", err)
		}
	case <-srv.Done:
		// 服务自行停止（多半是 Serve 出错），错误由下面的 Wait 返回。
	}
	return srv.Wait()
}
