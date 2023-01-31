package view

import (
	"context"
	"errors"

	"github.com/one2nc/cloud-lens/internal"
)

func extractApp(ctx context.Context) (*App, error) {
	app, ok := ctx.Value(internal.KeyApp).(*App)
	if !ok {
		return nil, errors.New("No application found in context")
	}

	return app, nil
}
