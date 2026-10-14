package internal

import (
	"context"

	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/gofs"
)

func Restarter(gs igs.Service) gofs.RestarterFunc {
	return func(ctx context.Context) error {
		if gs != nil {
			return gs.Restart()
		}
		return nil
	}
}
