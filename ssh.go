package leap0

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- SSH Types ---

type SSHAccess struct {
	ID         string `json:"id"`
	SandboxID  string `json:"sandbox_id"`
	Password   string `json:"password"`
	SSHCommand string `json:"ssh_command"`
	ExpiresAt  string `json:"expires_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SSHValidation struct {
	Valid bool `json:"valid"`
}

// --- SSHService ---

type SSHService struct{ t *transport.Client }

func (s *SSHService) CreateAccess(ctx context.Context, id string) (*SSHAccess, error) {
	var r SSHAccess
	return &r, wrapErr("create ssh", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/ssh/access", id), nil, &r, nil))
}

func (s *SSHService) DeleteAccess(ctx context.Context, id, credID string) error {
	return wrapErr("delete ssh", s.t.JSON(ctx, http.MethodDelete, fmt.Sprintf("/v1/sandbox/%s/ssh/%s", id, credID), nil, nil, nil))
}

func (s *SSHService) Validate(ctx context.Context, id, credID, password string) (*SSHValidation, error) {
	var r SSHValidation
	return &r, wrapErr("validate ssh", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/ssh/%s/validate", id, credID), map[string]any{"password": password}, &r, nil))
}

func (s *SSHService) Regenerate(ctx context.Context, id, credID string) (*SSHAccess, error) {
	var r SSHAccess
	return &r, wrapErr("regenerate ssh", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/ssh/%s/regen", id, credID), nil, &r, nil))
}
