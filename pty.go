package leap0

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- PTY Types ---

type PtySession struct {
	ID   string `json:"id"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
	Pid  int    `json:"pid,omitempty"`
}

type CreatePtyParams struct {
	Cols    int               `json:"cols,omitempty"`
	Rows    int               `json:"rows,omitempty"`
	Command string            `json:"command,omitempty"`
	Cwd     string            `json:"cwd,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

type ResizePtyParams struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

// --- PtyService ---

type PtyService struct{ t *transport.Client }

func (s *PtyService) ep(id, suffix string) string {
	return fmt.Sprintf("/v1/sandbox/%s/pty%s", id, suffix)
}

func (s *PtyService) List(ctx context.Context, id string) ([]PtySession, error) {
	var r []PtySession
	return r, wrapErr("list pty", s.t.JSON(ctx, http.MethodGet, s.ep(id, ""), nil, &r, nil))
}

func (s *PtyService) Create(ctx context.Context, id string, p *CreatePtyParams) (*PtySession, error) {
	var r PtySession
	return &r, wrapErr("create pty", s.t.JSON(ctx, http.MethodPost, s.ep(id, ""), p, &r, nil))
}

func (s *PtyService) Get(ctx context.Context, id, sessionID string) (*PtySession, error) {
	var r PtySession
	return &r, wrapErr("get pty", s.t.JSON(ctx, http.MethodGet, s.ep(id, "/"+sessionID), nil, &r, nil))
}

func (s *PtyService) Delete(ctx context.Context, id, sessionID string) error {
	return wrapErr("delete pty", s.t.JSON(ctx, http.MethodDelete, s.ep(id, "/"+sessionID), nil, nil, nil))
}

func (s *PtyService) Resize(ctx context.Context, id, sessionID string, p *ResizePtyParams) (*PtySession, error) {
	var r PtySession
	return &r, wrapErr("resize pty", s.t.JSON(ctx, http.MethodPost, s.ep(id, "/"+sessionID+"/resize"), p, &r, nil))
}
