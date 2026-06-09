package leap0

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

type TemplateImageConfig struct {
	Entrypoint []string          `json:"entrypoint,omitempty"`
	Cmd        []string          `json:"cmd,omitempty"`
	WorkingDir string            `json:"workingDir,omitempty"`
	User       string            `json:"user,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
}

type Template struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Digest      string               `json:"digest"`
	ImageConfig *TemplateImageConfig `json:"imageConfig,omitempty"`
	IsSystem    bool                 `json:"isSystem"`
	CreatedAt   time.Time            `json:"created_at"`
}

type CreateTemplateParams struct {
	Name        string              `json:"name"`
	URI         string              `json:"uri"`
	Credentials *RegistryCredential `json:"credentials,omitempty"`
}

type ListTemplatesParams struct {
	Query    string `json:"query,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type TemplateList struct {
	Items      []Template `json:"items"`
	TotalItems int        `json:"total_items"`
}

type RenameTemplateParams struct {
	Name string `json:"name"`
}

type RegistryCredential struct {
	Type                  string `json:"type"`
	Username              string `json:"username,omitempty"`
	Password              string `json:"password,omitempty"`
	AWSAccessKeyID        string `json:"aws_access_key_id,omitempty"`
	AWSSecretAccessKey    string `json:"aws_secret_access_key,omitempty"`
	AWSRegion             string `json:"aws_region,omitempty"`
	GCPServiceAccountJSON string `json:"gcp_service_account_json,omitempty"`
	AzureClientID         string `json:"azure_client_id,omitempty"`
	AzureClientSecret     string `json:"azure_client_secret,omitempty"`
	AzureTenantID         string `json:"azure_tenant_id,omitempty"`
}

type TemplatesService struct{ t *transport.Client }

func (s *TemplatesService) Create(ctx context.Context, p *CreateTemplateParams) (*Template, error) {
	var r Template
	return &r, wrapErr("create template", s.t.JSON(ctx, http.MethodPost, "/v1/template", p, &r, nil))
}

func (s *TemplatesService) List(ctx context.Context, p *ListTemplatesParams) (*TemplateList, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Page > 0 {
			opts.Query["page"] = strconv.Itoa(p.Page)
		}
		if p.PageSize > 0 {
			opts.Query["page_size"] = strconv.Itoa(p.PageSize)
		}
		if p.Query != "" {
			opts.Query["query"] = p.Query
		}
	}
	var r TemplateList
	return &r, wrapErr("list templates", s.t.JSON(ctx, http.MethodGet, "/v1/template", nil, &r, opts))
}

func (s *TemplatesService) Rename(ctx context.Context, id string, p *RenameTemplateParams) error {
	return wrapErr("rename template", s.t.JSON(ctx, http.MethodPatch, fmt.Sprintf("/v1/template/%s", id), p, nil, nil))
}

func (s *TemplatesService) Delete(ctx context.Context, id string) error {
	return wrapErr("delete template", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/template/%s", id), nil, nil, nil))
}
