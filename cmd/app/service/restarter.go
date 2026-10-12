package service

import (
	"context"

	"github.com/xxl6097/go-service/pkg/gs/igs"
)

type Restarter struct {
	gs igs.Service
}

func NewRestarter(gs igs.Service) *Restarter {
	return &Restarter{gs: gs}
}

func (re *Restarter) Restart(ctx context.Context) error {
	if re.gs != nil {
		return re.gs.Restart()
	}
	return nil
}
