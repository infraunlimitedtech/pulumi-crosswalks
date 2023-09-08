package manager

import (
	"automation-api/common/apiserver"
	"automation-api/common/log"
	"automation-api/common/pulumi"
	"automation-api/hetzner-snapshots-manager/hetzner"
	"context"
	"fmt"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"go.uber.org/zap"
)

type Manager struct {
	ctx       context.Context
	Runner    *pulumi.Pulumi
	APIServer *apiserver.Server
	Hetzner   *hetzner.API
	Logger    *zap.Logger
}

func New(ctx context.Context) (*Manager, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}

	logger, _ := log.New(config.Debug)

	runner, err := pulumi.New(ctx, logger, config.Stack.Name, config.Stack.Path)
	if err != nil {
		return nil, fmt.Errorf("create pulumi runner: %w", err)
	}

	token, err := runner.GetConfig("hcloud:token")
	if err != nil {
		return nil, fmt.Errorf("retrieve hcloud token: %w", err)
	}

	hetzner := hetzner.New(ctx, logger, token)

	snapshots, err := hetzner.GatherSnapshotInfo()
	if err != nil {
		return nil, fmt.Errorf("retrieve info about snapshots: %w", err)
	}

	httpAddr := fmt.Sprintf("localhost:%d", config.APIServerPort)

	apiServer, err := apiserver.New(httpAddr, logger, getAllRoutes(snapshots))
	if err != nil {
		return nil, fmt.Errorf("create api server: %w", err)
	}

	// pass needed info to pulumi cli. Bad naming, but I like to `attach` :)
	runner.AttachToAPIServer(apiServer.Addr())

	return &Manager{
		ctx:       ctx,
		Runner:    runner,
		APIServer: apiServer,
		Logger:    logger,
		Hetzner:   hetzner,
	}, nil
}

func (m *Manager) ProcessEvents(events []events.EngineEvent, dryRun bool) error {
	for _, p := range events {
		if p.ResourcePreEvent != nil {
			metadata := p.ResourcePreEvent.Metadata
			if metadata.Op == "delete" && metadata.Type == "hcloud:index/server:Server" {
				m.Logger.Info(fmt.Sprintf("we are deleting a resource %s/%s. a snapshot creation needed",
					metadata.Old.Type, metadata.Old.ID,
				))

				if dryRun {
					m.Logger.Info("dry run, skipping snapshot creation")
					continue
				}

				err := m.makeSnapshot(metadata.Old.ID)
				if err != nil {
					return fmt.Errorf("make snapshot: %w", err)
				}
			}
		}
	}

	return nil
}

func (m *Manager) makeSnapshot(id string) error {
	timeout := 20 * time.Minute
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	m.Logger.Info(fmt.Sprintf("the snapshot creation for %s may take a time. Please be patient. Max allowed time is %s",
		id, timeout.String(),
	))

	start := time.Now()
	err := m.Hetzner.CreateSnapshot(ctx, id)

	m.Logger.Debug(fmt.Sprintf("the snapshot creation took %s", time.Since(start).String()))

	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	return nil
}
