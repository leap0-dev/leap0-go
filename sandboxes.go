package leap0

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

type SandboxState string

const (
	StateStarting     SandboxState = "starting"
	StateRunning      SandboxState = "running"
	StatePaused       SandboxState = "paused"
	StateUnpausing    SandboxState = "unpausing"
	StateSnapshotting SandboxState = "snapshotting"
	StateDeleting     SandboxState = "deleting"
	StateStopped      SandboxState = "stopped"
	StateStopping     SandboxState = "stopping"
	StateDeleted      SandboxState = "deleted"
)

type NetworkPolicyMode string

const (
	NetworkPolicyAllowAll NetworkPolicyMode = "allow-all"
	NetworkPolicyDenyAll  NetworkPolicyMode = "deny-all"
	NetworkPolicyCustom   NetworkPolicyMode = "custom"
)

type SandboxData struct {
	ID            string            `json:"id"`
	TemplateID    string            `json:"template_id,omitempty"`
	TemplateName  string            `json:"template_name,omitempty"`
	State         SandboxState      `json:"state"`
	VCPU          int               `json:"vcpu"`
	Memory        int               `json:"memory"`
	Disk          int               `json:"disk,omitempty"`
	Timeout       int               `json:"timeout"`
	AutoPause     bool              `json:"auto_pause"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
	NetworkPolicy *NetworkPolicy    `json:"network_policy,omitempty"`
	Mounts        []Mount           `json:"mounts,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
}

type NetworkPolicy struct {
	Mode         NetworkPolicyMode  `json:"mode"`
	AllowDomains []string           `json:"allow_domains,omitempty"`
	AllowCIDRs   []string           `json:"allow_cidrs,omitempty"`
	Transforms   []NetworkTransform `json:"transforms,omitempty"`
}

type NetworkTransform struct {
	Domain        string            `json:"domain"`
	InjectHeaders map[string]string `json:"inject_headers,omitempty"`
	StripHeaders  []string          `json:"strip_headers,omitempty"`
}

type Mount struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Bucket    string `json:"bucket"`
	MountPath string `json:"mount_path"`
	Prefix    string `json:"prefix,omitempty"`
	ReadOnly  bool   `json:"read_only"`
}
type MountRequest struct {
	Bucket          string `json:"bucket"`
	MountPath       string `json:"mount_path"`
	Endpoint        string `json:"endpoint"`
	Prefix          string `json:"prefix,omitempty"`
	ReadOnly        bool   `json:"read_only"`
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
}

type CreateSandboxParams struct {
	TemplateName  string            `json:"template_name"`
	VCPU          int               `json:"vcpu,omitempty"`
	Memory        int               `json:"memory,omitempty"`
	Timeout       int               `json:"timeout,omitempty"`
	AutoPause     *bool             `json:"auto_pause,omitempty"`
	OtelExport    *bool             `json:"otel_export,omitempty"`
	Telemetry     *bool             `json:"telemetry,omitempty"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
	NetworkPolicy *NetworkPolicy    `json:"network_policy,omitempty"`
	Mounts        []MountRequest    `json:"mounts,omitempty"`
}

type ListSandboxesParams struct {
	Page     int
	PageSize int
	State    string
	Sort     string
	OrderBy  string
}

type SandboxListItem struct {
	ID              string       `json:"id"`
	TemplateID      string       `json:"template_id"`
	State           SandboxState `json:"state"`
	LaunchTime      string       `json:"launch_time,omitempty"`
	StateChangeTime string       `json:"state_change_time,omitempty"`
	TimeoutAt       int          `json:"timeout_at,omitempty"`
	CreatedAt       string       `json:"created_at"`
}

type SandboxList struct {
	Items      []SandboxData `json:"items"`
	TotalItems int           `json:"total_items"`
}

type PresignedURLParams struct {
	Port      int `json:"port"`
	ExpiresIn int `json:"expires_in,omitempty"`
}

type PresignedURL struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	URL       string    `json:"url"`
	SandboxID string    `json:"sandbox_id"`
	Port      int       `json:"port"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt string    `json:"created_at"`
}

type SandboxesService struct{ t *transport.Client }

func (s *SandboxesService) Create(ctx context.Context, p *CreateSandboxParams) (*SandboxData, error) {
	if p == nil {
		p = &CreateSandboxParams{}
	}
	var r SandboxData
	return &r, wrapErr("create sandbox", s.t.JSON(ctx, http.MethodPost, "/v1/sandbox", p, &r, nil))
}

func (s *SandboxesService) Get(ctx context.Context, id string) (*SandboxData, error) {
	var r SandboxData
	return &r, wrapErr("get sandbox", s.t.JSON(ctx, http.MethodGet, fmt.Sprintf("/v1/sandbox/%s", id), nil, &r, nil))
}

func (s *SandboxesService) List(ctx context.Context, p *ListSandboxesParams) (*SandboxList, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Page > 0 {
			opts.Query["page"] = strconv.Itoa(p.Page)
		}
		if p.PageSize > 0 {
			opts.Query["page-size"] = strconv.Itoa(p.PageSize)
		}
		if p.State != "" {
			opts.Query["state"] = p.State
		}
		if p.Sort != "" {
			opts.Query["sort"] = p.Sort
		}
		if p.OrderBy != "" {
			opts.Query["order-by"] = p.OrderBy
		}
	}
	var r SandboxList
	return &r, wrapErr("list sandboxes", s.t.JSON(ctx, http.MethodGet, "/v1/sandboxes", nil, &r, opts))
}

func (s *SandboxesService) Pause(ctx context.Context, id string) error {
	return wrapErr("pause", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/pause", id), nil, nil, nil))
}

func (s *SandboxesService) Stop(ctx context.Context, id string) error {
	return wrapErr("stop", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/stop", id), nil, nil, nil))
}

func (s *SandboxesService) Start(ctx context.Context, id string) (*SandboxData, error) {
	var r SandboxData
	return &r, wrapErr("start", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/start", id), nil, &r, nil))
}

func (s *SandboxesService) Delete(ctx context.Context, id string) error {
	return wrapErr("delete", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/sandbox/%s", id), nil, nil, nil))
}

func (s *SandboxesService) CreateSnapshot(ctx context.Context, id string, p *CreateSnapshotParams) (*Snapshot, error) {
	if p == nil {
		p = &CreateSnapshotParams{}
	}
	var r Snapshot
	return &r, wrapErr("create snapshot", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/snapshot/create", id), p, &r, nil))
}

func (s *SandboxesService) AddMount(ctx context.Context, id string, m *MountRequest) error {
	return wrapErr("add mount", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/mounts", id), m, nil, nil))
}

func (s *SandboxesService) UpdateMount(ctx context.Context, id, mountID string, m *MountRequest) error {
	return wrapErr("update mount", s.t.JSON(ctx, http.MethodPatch, fmt.Sprintf("/v1/sandbox/%s/mounts/%s", id, mountID), m, nil, nil))
}

func (s *SandboxesService) DeleteMount(ctx context.Context, id, mountID string) error {
	return wrapErr("delete mount", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/sandbox/%s/mounts/%s", id, mountID), nil, nil, nil))
}

func (s *SandboxesService) GetWorkdir(ctx context.Context, id string) (string, error) {
	var r struct {
		Workdir string `json:"workdir"`
	}
	err := s.t.JSON(ctx, http.MethodGet, fmt.Sprintf("/v1/sandbox/%s/system/workdir", id), nil, &r, nil)
	return r.Workdir, wrapErr("get workdir", err)
}

func (s *SandboxesService) GetHomeDir(ctx context.Context, id string) (string, error) {
	var r struct {
		HomeDir string `json:"home_dir"`
	}
	err := s.t.JSON(ctx, http.MethodGet, fmt.Sprintf("/v1/sandbox/%s/system/user-home-dir", id), nil, &r, nil)
	return r.HomeDir, wrapErr("get home dir", err)
}

func (s *SandboxesService) CreatePresignedURL(ctx context.Context, id string, p *PresignedURLParams) (*PresignedURL, error) {
	var r PresignedURL
	return &r, wrapErr("create presigned url", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/presigned-url", id), p, &r, nil))
}

func (s *SandboxesService) DeletePresignedURL(ctx context.Context, id, urlID string) error {
	return wrapErr("delete presigned url", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/sandbox/%s/presigned-url/%s", id, urlID), nil, nil, nil))
}

func (s *SandboxesService) InvokeURL(id, path string, port int) string {
	return s.t.SandboxURL(id, port, path)
}

func (s *SandboxesService) WebsocketURL(id, path string, port int) string {
	return s.t.SandboxWSURL(id, port, path)
}
