package leap0

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

type Snapshot struct {
	ID            string         `json:"id"`
	Name          string         `json:"name,omitempty"`
	TemplateID    string         `json:"template_id,omitempty"`
	VCPU          int            `json:"vcpu,omitempty"`
	Memory        int            `json:"memory,omitempty"`
	Disk          int            `json:"disk,omitempty"`
	State         SandboxState   `json:"state,omitempty"`
	NetworkPolicy *NetworkPolicy `json:"network_policy,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}

type CreateSnapshotParams struct {
	Name             string `json:"name,omitempty"`
	KillSandboxAfter bool   `json:"kill_sandbox_after,omitempty"`
}

type ListSnapshotsParams struct {
	Query    string `json:"query,omitempty"`
	Sort     string
	OrderBy  string
	Page     int
	PageSize int
}

type SnapshotList struct {
	Items      []Snapshot `json:"items"`
	TotalItems int        `json:"total_items"`
}

type RestoreSnapshotParams struct {
	SnapshotName  string            `json:"snapshot_name"`
	Timeout       int               `json:"timeout,omitempty"`
	AutoPause     *bool             `json:"auto_pause,omitempty"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
	NetworkPolicy *NetworkPolicy    `json:"network_policy,omitempty"`
}

type SnapshotsService struct{ t *transport.Client }

func (s *SnapshotsService) List(ctx context.Context, p *ListSnapshotsParams) (*SnapshotList, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Page > 0 {
			opts.Query["page"] = strconv.Itoa(p.Page)
		}
		if p.PageSize > 0 {
			opts.Query["page-size"] = strconv.Itoa(p.PageSize)
		}
		if p.Query != "" {
			opts.Query["query"] = p.Query
		}
		if p.Sort != "" {
			opts.Query["sort"] = p.Sort
		}
		if p.OrderBy != "" {
			opts.Query["order-by"] = p.OrderBy
		}
	}
	var r SnapshotList
	return &r, wrapErr("list snapshots", s.t.JSON(ctx, http.MethodGet, "/v1/snapshots", nil, &r, opts))
}

func (s *SnapshotsService) Restore(ctx context.Context, p *RestoreSnapshotParams) (*SandboxData, error) {
	var r SandboxData
	return &r, wrapErr("restore snapshot", s.t.JSON(ctx, http.MethodPost, "/v1/snapshot/restore", p, &r, nil))
}

func (s *SnapshotsService) Delete(ctx context.Context, id string) error {
	return wrapErr("delete snapshot", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/snapshot/%s", id), nil, nil, nil))
}
