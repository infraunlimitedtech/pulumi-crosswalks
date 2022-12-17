package hetzner

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	managedByLabel  = "managed-by"
	timestampLabel  = "timestamp"
	serverNameLabel = "server-name"
	app             = "pulumi-automation-api"
)

var ErrSnapshotNotFound = errors.New("snapshot not found")

type SnapshotInfo struct {
	ID         int
	ServerName string `json:"server_name,omitempty"`
	Timestamp  string `json:",omitempty"`
}

type Snapshots struct {
	gathered []SnapshotInfo
}

func (h *API) GatherSnapshotInfo() (*Snapshots, error) {
	var gathered []SnapshotInfo

	opts := hcloud.ImageListOpts{
		Type: []hcloud.ImageType{hcloud.ImageTypeSnapshot},
	}

	images, err := h.client.Image.AllWithOpts(h.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("gather snapshots info: %w", err)
	}

	for _, image := range images {
		if image.Labels[managedByLabel] != "" && image.Labels[managedByLabel] == app {
			if image.Labels[timestampLabel] != "" {
				if image.Labels[serverNameLabel] != "" {
					i := SnapshotInfo{
						ID:         image.ID,
						ServerName: image.Labels[serverNameLabel],
						Timestamp:  image.Labels[timestampLabel],
					}
					h.logger.Debug(fmt.Sprintf("add snapshot %+v to list available images", i))
					gathered = append(gathered, i)
				}
			}
		}
	}

	if len(gathered) == 0 {
		h.logger.Info("No available images found in Hetzner cloud")
	}

	return &Snapshots{
		gathered: gathered,
	}, nil
}

func (h *API) CreateSnapshot(ctx context.Context, idOrName string) error {
	srv, err := h.getServer(ctx, idOrName)
	if err != nil {
		return fmt.Errorf("get a server: %w", err)
	}

	labels := map[string]string{
		serverNameLabel: srv.Name,
		managedByLabel:  app,
		timestampLabel:  strconv.FormatInt(time.Now().Unix(), 10),
	}

	description := fmt.Sprintf("automatically made for %s at %s",
		srv.Name,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	opts := &hcloud.ServerCreateImageOpts{
		Description: &description,
		Type:        hcloud.ImageTypeSnapshot,
		Labels:      labels,
	}

	snapshot, _, err := h.client.Server.CreateImage(ctx, srv, opts)
	if err != nil {
		return fmt.Errorf("create a image from server: %w", err)
	}

	for {
		resp, _, err := h.client.Image.GetByID(ctx, snapshot.Image.ID)
		if err != nil {
			return fmt.Errorf("get the created snapshot: %w", err)
		}

		if resp.Status != hcloud.ImageStatusAvailable {
			time.Sleep(1 * time.Second)

			continue
		}

		h.logger.Info(fmt.Sprintf("snapshot %d is ready", resp.ID))

		break
	}

	return nil
}

func (s *Snapshots) GetLastForServer(target string) (*SnapshotInfo, error) {
	var allByTarget []SnapshotInfo
	allByTarget, err := s.getAllByName(target)
	if err != nil {
		return nil, fmt.Errorf("get last snapshot: %w", err)
	}

	sort.Slice(allByTarget, func(i, j int) bool {
		return allByTarget[i].Timestamp > allByTarget[j].Timestamp
	})

	return &allByTarget[0], nil
}

func (s *Snapshots) getAllByName(target string) ([]SnapshotInfo, error) {
	var all []SnapshotInfo

	for _, snapshot := range s.gathered {
		if snapshot.ServerName == target {
			all = append(all, snapshot)
		}
	}
	if len(all) == 0 {
		return all, ErrSnapshotNotFound
	}

	return all, nil
}
