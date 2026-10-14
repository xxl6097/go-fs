package internal

import (
	"context"

	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/gofs"
)

func Uninstaller(gs igs.Service) gofs.UninstallFunc {
	return func(ctx context.Context) error {
		if gs != nil {
			return gs.UnInstall()
		}
		return nil
	}
}
