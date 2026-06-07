package leap0

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Sandbox Types ---

type SandboxState string

const (
	StateStarting     SandboxState = "starting"
	StateRunning      SandboxState = "running"
	StatePaused       SandboxState = "paused"
	StateUnpausing    SandboxState = "unpausing"
	StateSnapshotting SandboxState = "snapshotting"
	StateDeleting     SandboxState = "deleting"
	StateStopped      SandboxState = "stopped"
)

type SandboxData struct {
	ID            string            `json:"id"`
	Name          string            `json:"name,omitempty"`
	State         SandboxState      `json:"state"`
	TemplateName  string            `json:"template_name"`
	VCPU          int               `json:"vcpu"`
	Memory        int               `json:"memory"`
	Timeout       int               `json:"timeout"`
	AutoPause     bool              `json:"auto_pause"`
	OtelExport    bool              `json:"otel_export"`
	Telemetry     bool              `json:"telemetry"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
	NetworkPolicy *NetworkPolicy    `json:"network_policy,omitempty"`
	Mounts        []Mount           `json:"mounts,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type NetworkPolicy struct {
	AllowOutbound bool          `json:"allow_outbound"`
	AllowInbound  bool          `json:"allow_inbound"`
	AllowedHosts  []AllowedHost `json:"allowed_hosts,omitempty"`
}

type AllowedHost struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports,omitempty"`
}

type Mount struct {
	ID         string `json:"id,omitempty"`
	MountPath  string `json:"mount_path"`
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
	Provider   string `json:"provider"`
	AccessKey  string `json:"access_key,omitempty"`
	SecretKey  string `json:"secret_key,omitempty"`
	ReadOnly   bool   `json:"read_only"`
}

type CreateSandboxParams struct {
	TemplateName  string            `json:"template_name,omitempty"`
	VCPU          int               `json:"vcpu,omitempty"`
	Memory        int               `json:"memory,omitempty"`
	Timeout       int               `json:"timeout,omitempty"`
	AutoPause     *bool             `json:"auto_pause,omitempty"`
	OtelExport    *bool             `json:"otel_export,omitempty"`
	Telemetry     *bool             `json:"telemetry,omitempty"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
	NetworkPolicy *NetworkPolicy    `json:"network_policy,omitempty"`
	Mounts        []Mount           `json:"mounts,omitempty"`
}

type ListSandboxesParams struct {
	Page     int
	PageSize int
	State    string
	Sort     string
	OrderBy  string
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
	URL       string    `json:"url"`
	Port      int       `json:"port"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Snapshot struct {
	ID           string    `json:"id"`
	SandboxID    string    `json:"sandbox_id"`
	TemplateName string    `json:"template_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListSnapshotsParams struct {
	Page      int
	PageSize  int
	SandboxID string
}

type SnapshotList struct {
	Items      []Snapshot `json:"items"`
	TotalItems int        `json:"total_items"`
}

type RestoreSnapshotParams struct {
	SnapshotID string            `json:"snapshot_id"`
	VCPU       int               `json:"vcpu,omitempty"`
	Memory     int               `json:"memory,omitempty"`
	Timeout    int               `json:"timeout,omitempty"`
	AutoPause  *bool             `json:"auto_pause,omitempty"`
	EnvVars    map[string]string `json:"env_vars,omitempty"`
}

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTemplateParams struct {
	Name       string              `json:"name"`
	Dockerfile string              `json:"dockerfile"`
	Credential *RegistryCredential `json:"credential,omitempty"`
}

type RenameTemplateParams struct {
	Name string `json:"name"`
}

type RegistryCredential struct {
	Type              string `json:"type"`
	Server            string `json:"server,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	Region            string `json:"region,omitempty"`
	AccessKeyID       string `json:"access_key_id,omitempty"`
	SecretAccessKey    string `json:"secret_access_key,omitempty"`
	ServiceAccountKey string `json:"service_account_key,omitempty"`
	TenantID          string `json:"tenant_id,omitempty"`
	ClientID          string `json:"client_id,omitempty"`
	ClientSecret      string `json:"client_secret,omitempty"`
}

// --- SandboxesService ---

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
	return &r, wrapErr("get sandbox", s.t.JSON(ctx, http.MethodGet, fmt.Sprintf("/v1/sandbox/%s/", id), nil, &r, nil))
}

func (s *SandboxesService) List(ctx context.Context, p *ListSandboxesParams) (*SandboxList, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Page > 0 {
			opts.Query["page"] = strconv.Itoa(p.Page)
		}
		if p.PageSize > 0 {
			opts.Query["page_size"] = strconv.Itoa(p.PageSize)
		}
		if p.State != "" {
			opts.Query["state"] = p.State
		}
		if p.Sort != "" {
			opts.Query["sort"] = p.Sort
		}
		if p.OrderBy != "" {
			opts.Query["order_by"] = p.OrderBy
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

func (s *SandboxesService) Start(ctx context.Context, id string) error {
	return wrapErr("start", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/start", id), nil, nil, nil))
}

func (s *SandboxesService) Delete(ctx context.Context, id string) error {
	return wrapErr("delete", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/sandbox/%s/", id), nil, nil, nil))
}

func (s *SandboxesService) CreateSnapshot(ctx context.Context, id string) (*Snapshot, error) {
	var r Snapshot
	return &r, wrapErr("create snapshot", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/snapshot/create", id), nil, &r, nil))
}

func (s *SandboxesService) AddMount(ctx context.Context, id string, m *Mount) error {
	return wrapErr("add mount", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/mounts", id), m, nil, nil))
}

func (s *SandboxesService) UpdateMount(ctx context.Context, id, mountID string, m *Mount) error {
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

// --- SnapshotsService ---

type SnapshotsService struct{ t *transport.Client }

func (s *SnapshotsService) List(ctx context.Context, p *ListSnapshotsParams) (*SnapshotList, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Page > 0 {
			opts.Query["page"] = strconv.Itoa(p.Page)
		}
		if p.PageSize > 0 {
			opts.Query["page_size"] = strconv.Itoa(p.PageSize)
		}
		if p.SandboxID != "" {
			opts.Query["sandbox_id"] = p.SandboxID
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

// --- TemplatesService ---

type TemplatesService struct{ t *transport.Client }

func (s *TemplatesService) Create(ctx context.Context, p *CreateTemplateParams) (*Template, error) {
	var r Template
	return &r, wrapErr("create template", s.t.JSON(ctx, http.MethodPost, "/v1/template", p, &r, nil))
}

func (s *TemplatesService) Rename(ctx context.Context, id string, p *RenameTemplateParams) error {
	return wrapErr("rename template", s.t.JSON(ctx, http.MethodPatch, fmt.Sprintf("/v1/template/%s", id), p, nil, nil))
}

func (s *TemplatesService) Delete(ctx context.Context, id string) error {
	return wrapErr("delete template", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/template/%s", id), nil, nil, nil))
}
