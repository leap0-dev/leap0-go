package leap0

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Process Types ---

type ExecParams struct {
	Command string            `json:"command"`
	Cwd     string            `json:"cwd,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Timeout int               `json:"timeout,omitempty"`
}

type ExecResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

// --- ProcessService ---

type ProcessService struct{ t *transport.Client }

func (s *ProcessService) Execute(ctx context.Context, id string, p *ExecParams) (*ExecResult, error) {
	var r ExecResult
	return &r, wrapErr("execute", s.t.JSON(ctx, http.MethodPost, fmt.Sprintf("/v1/sandbox/%s/process/execute", id), p, &r, nil))
}
