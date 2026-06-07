package leap0

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Code Interpreter Types ---

type Language string

const (
	Python     Language = "python"
	TypeScript Language = "typescript"
)

type CodeContext struct {
	ID       string   `json:"id"`
	Language Language `json:"language"`
	Cwd      string   `json:"cwd,omitempty"`
}

type RunCodeParams struct {
	Code      string            `json:"code"`
	Language  Language          `json:"language,omitempty"`
	ContextID string            `json:"context_id,omitempty"`
	EnvVars   map[string]string `json:"env_vars,omitempty"`
	TimeoutMs int               `json:"timeout_ms,omitempty"`
}

type RunCodeResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

type StreamEventKind string

const (
	EventStdout StreamEventKind = "stdout"
	EventStderr StreamEventKind = "stderr"
	EventExit   StreamEventKind = "exit"
	EventError  StreamEventKind = "error"
)

type StreamEvent struct {
	Type     StreamEventKind `json:"type"`
	Data     string          `json:"data,omitempty"`
	ExitCode int             `json:"exit_code,omitempty"`
}

// --- CodeInterpreterService ---

type CodeInterpreterService struct{ t *transport.Client }

func (s *CodeInterpreterService) url(id, path string) string {
	return s.t.SandboxURL(id, 0, path)
}

func (s *CodeInterpreterService) Health(ctx context.Context, id string) (bool, error) {
	resp, err := s.t.Raw(ctx, http.MethodGet, s.url(id, "/healthz"), nil, nil)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return resp.StatusCode == 200, nil
}

func (s *CodeInterpreterService) CreateContext(ctx context.Context, id string, lang Language, cwd string) (*CodeContext, error) {
	body := map[string]any{}
	if lang != "" {
		body["language"] = lang
	}
	if cwd != "" {
		body["cwd"] = cwd
	}
	var r CodeContext
	return &r, wrapErr("create context", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/contexts"), body, &r, nil))
}

func (s *CodeInterpreterService) ListContexts(ctx context.Context, id string) ([]CodeContext, error) {
	var r []CodeContext
	return r, wrapErr("list contexts", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/contexts"), nil, &r, nil))
}

func (s *CodeInterpreterService) GetContext(ctx context.Context, id, contextID string) (*CodeContext, error) {
	var r CodeContext
	return &r, wrapErr("get context", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/contexts/"+contextID), nil, &r, nil))
}

func (s *CodeInterpreterService) DeleteContext(ctx context.Context, id, contextID string) error {
	return wrapErr("delete context", s.t.JSONAbsolute(ctx, http.MethodDelete, s.url(id, "/contexts/"+contextID), nil, nil, nil))
}

func (s *CodeInterpreterService) Run(ctx context.Context, id string, p *RunCodeParams) (*RunCodeResult, error) {
	var r RunCodeResult
	return &r, wrapErr("run code", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/execute"), p, &r, nil))
}

func (s *CodeInterpreterService) RunStream(ctx context.Context, id string, p *RunCodeParams) (<-chan StreamEvent, <-chan error, error) {
	stream, err := s.t.SSEAbsolute(ctx, http.MethodPost, s.url(id, "/execute/async"), p, nil)
	if err != nil {
		return nil, nil, wrapErr("run stream", err)
	}

	events := make(chan StreamEvent)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)
		defer stream.Close()

		for {
			raw, err := stream.Next()
			if err != nil {
				if err != io.EOF {
					errs <- err
				}
				return
			}
			var ev StreamEvent
			if err := json.Unmarshal(raw, &ev); err != nil {
				errs <- fmt.Errorf("unmarshal event: %w", err)
				return
			}
			select {
			case events <- ev:
			case <-ctx.Done():
				return
			}
		}
	}()

	return events, errs, nil
}
