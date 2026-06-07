package leap0

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Desktop Types ---

type DisplayInfo struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Window struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ClickParams struct {
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Button string `json:"button,omitempty"`
	Count  int    `json:"count,omitempty"`
}

type DragParams struct {
	StartX int `json:"start_x"`
	StartY int `json:"start_y"`
	EndX   int `json:"end_x"`
	EndY   int `json:"end_y"`
}

type ScrollParams struct {
	X         int `json:"x,omitempty"`
	Y         int `json:"y,omitempty"`
	Direction int `json:"direction"`
}

type ScreenshotParams struct {
	Format  string `json:"format,omitempty"`
	Quality int    `json:"quality,omitempty"`
}

type RegionScreenshotParams struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Format  string `json:"format,omitempty"`
	Quality int    `json:"quality,omitempty"`
}

type RecordingStatus struct {
	Recording bool   `json:"recording"`
	ID        string `json:"id,omitempty"`
}

type Recording struct {
	ID        string `json:"id"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
}

type DesktopHealth struct {
	Status string `json:"status"`
}

type ProcessInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Pid    int    `json:"pid,omitempty"`
}

type ProcessList struct {
	Processes []ProcessInfo `json:"processes"`
}

type ProcessLogs struct {
	Logs string `json:"logs"`
}

type StatusEvent struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

// --- DesktopService ---

type DesktopService struct{ t *transport.Client }

func (s *DesktopService) url(id, path string) string {
	return s.t.SandboxURL(id, 0, path)
}

func (s *DesktopService) URL(id string) string {
	return s.url(id, "")
}

func (s *DesktopService) DisplayInfo(ctx context.Context, id string) (*DisplayInfo, error) {
	var r DisplayInfo
	return &r, wrapErr("display info", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/display"), nil, &r, nil))
}

func (s *DesktopService) Screen(ctx context.Context, id string) (*DisplayInfo, error) {
	var r DisplayInfo
	return &r, wrapErr("screen", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/display/screen"), nil, &r, nil))
}

func (s *DesktopService) ResizeScreen(ctx context.Context, id string, width, height int) (*DisplayInfo, error) {
	var r DisplayInfo
	return &r, wrapErr("resize", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/display/screen"), map[string]any{"width": width, "height": height}, &r, nil))
}

func (s *DesktopService) Windows(ctx context.Context, id string) ([]Window, error) {
	var r []Window
	return r, wrapErr("windows", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/display/windows"), nil, &r, nil))
}

func (s *DesktopService) Screenshot(ctx context.Context, id string, p *ScreenshotParams) ([]byte, error) {
	opts := &transport.Options{Query: map[string]string{}}
	if p != nil {
		if p.Format != "" {
			opts.Query["format"] = p.Format
		}
		if p.Quality > 0 {
			opts.Query["quality"] = fmt.Sprintf("%d", p.Quality)
		}
	}
	return s.t.BytesAbsolute(ctx, http.MethodGet, s.url(id, "/api/screenshot"), nil, opts)
}

func (s *DesktopService) ScreenshotRegion(ctx context.Context, id string, p *RegionScreenshotParams) ([]byte, error) {
	return s.t.BytesAbsolute(ctx, http.MethodPost, s.url(id, "/api/screenshot/region"), p, nil)
}

func (s *DesktopService) PointerPosition(ctx context.Context, id string) (*Point, error) {
	var r Point
	return &r, wrapErr("pointer", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/input/position"), nil, &r, nil))
}

func (s *DesktopService) MovePointer(ctx context.Context, id string, x, y int) (*Point, error) {
	var r Point
	return &r, wrapErr("move", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/move"), map[string]any{"x": x, "y": y}, &r, nil))
}

func (s *DesktopService) Click(ctx context.Context, id string, p *ClickParams) (*Point, error) {
	var r Point
	return &r, wrapErr("click", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/click"), p, &r, nil))
}

func (s *DesktopService) Drag(ctx context.Context, id string, p *DragParams) (*Point, error) {
	var r Point
	return &r, wrapErr("drag", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/drag"), p, &r, nil))
}

func (s *DesktopService) Scroll(ctx context.Context, id string, p *ScrollParams) (*Point, error) {
	var r Point
	return &r, wrapErr("scroll", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/scroll"), p, &r, nil))
}

func (s *DesktopService) Type(ctx context.Context, id, text string) error {
	return wrapErr("type", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/type"), map[string]any{"text": text}, nil, nil))
}

func (s *DesktopService) PressKey(ctx context.Context, id, key string) error {
	return wrapErr("press", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/press"), map[string]any{"key": key}, nil, nil))
}

func (s *DesktopService) Hotkey(ctx context.Context, id string, keys []string) error {
	return wrapErr("hotkey", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/input/hotkey"), map[string]any{"keys": keys}, nil, nil))
}

func (s *DesktopService) StartRecording(ctx context.Context, id string) (*RecordingStatus, error) {
	var r RecordingStatus
	return &r, wrapErr("start recording", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/recording/start"), nil, &r, nil))
}

func (s *DesktopService) StopRecording(ctx context.Context, id string) (*RecordingStatus, error) {
	var r RecordingStatus
	return &r, wrapErr("stop recording", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/recording/stop"), nil, &r, nil))
}

func (s *DesktopService) Recordings(ctx context.Context, id string) ([]Recording, error) {
	var r []Recording
	return r, wrapErr("recordings", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/recordings"), nil, &r, nil))
}

func (s *DesktopService) DownloadRecording(ctx context.Context, id, recID string) ([]byte, error) {
	return s.t.BytesAbsolute(ctx, http.MethodGet, s.url(id, "/api/recordings/"+recID+"/download"), nil, nil)
}

func (s *DesktopService) DeleteRecording(ctx context.Context, id, recID string) error {
	return wrapErr("delete recording", s.t.JSONAbsolute(ctx, http.MethodDelete, s.url(id, "/api/recordings/"+recID), nil, nil, nil))
}

func (s *DesktopService) Health(ctx context.Context, id string) (*DesktopHealth, error) {
	var r DesktopHealth
	return &r, wrapErr("health", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/healthz"), nil, &r, nil))
}

func (s *DesktopService) ProcessStatus(ctx context.Context, id string) (*ProcessList, error) {
	var r ProcessList
	return &r, wrapErr("process status", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/status"), nil, &r, nil))
}

func (s *DesktopService) RestartProcess(ctx context.Context, id, name string) error {
	return wrapErr("restart", s.t.JSONAbsolute(ctx, http.MethodPost, s.url(id, "/api/process/"+name+"/restart"), nil, nil, nil))
}

func (s *DesktopService) ProcessLogs(ctx context.Context, id, name string) (*ProcessLogs, error) {
	var r ProcessLogs
	return &r, wrapErr("logs", s.t.JSONAbsolute(ctx, http.MethodGet, s.url(id, "/api/process/"+name+"/logs"), nil, &r, nil))
}

func (s *DesktopService) StatusStream(ctx context.Context, id string) (<-chan StatusEvent, <-chan error, error) {
	stream, err := s.t.SSEAbsolute(ctx, http.MethodGet, s.url(id, "/api/status/stream"), nil, nil)
	if err != nil {
		return nil, nil, wrapErr("status stream", err)
	}

	events := make(chan StatusEvent)
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
			var ev StatusEvent
			if err := json.Unmarshal(raw, &ev); err != nil {
				errs <- fmt.Errorf("unmarshal: %w", err)
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

func (s *DesktopService) WaitUntilReady(ctx context.Context, id string, timeout time.Duration) error {
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		h, err := s.Health(ctx, id)
		if err == nil && h.Status == "ok" {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return fmt.Errorf("desktop: not ready within %s: %w", timeout, ErrTimeout)
}
