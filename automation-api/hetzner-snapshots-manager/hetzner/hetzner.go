package hetzner

import (
	"context"
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"go.uber.org/zap"
)

var ErrServerNotFound = errors.New("server not found")

type API struct {
	ctx    context.Context
	logger *zap.Logger
	client *hcloud.Client
}

func New(ctx context.Context, logger *zap.Logger, token string) *API {
	return &API{
		client: hcloud.NewClient(hcloud.WithToken(token)),
		ctx:    ctx,
		logger: logger,
	}
}

func (h *API) getServer(ctx context.Context, idOrName string) (*hcloud.Server, error) {
	server, _, err := h.client.Server.Get(ctx, idOrName)

	if server == nil {
		return nil, ErrServerNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get server: %w", err)
	}

	return server, nil
}
