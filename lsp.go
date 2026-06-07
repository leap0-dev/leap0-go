package leap0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- LSP Types ---

type LspStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type LspResult struct {
	ID     int             `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *LspError       `json:"error,omitempty"`
}

type LspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LspServerParams struct {
	LanguageID    string `json:"language_id"`
	PathToProject string `json:"path_to_project"`
}

type LspDocumentParams struct {
	LanguageID    string `json:"language_id"`
	PathToProject string `json:"path_to_project"`
	URI           string `json:"uri"`
	Text          string `json:"text,omitempty"`
	Version       int    `json:"version,omitempty"`
}

type LspPositionParams struct {
	LanguageID    string `json:"language_id"`
	PathToProject string `json:"path_to_project"`
	URI           string `json:"uri"`
	Line          int    `json:"line"`
	Character     int    `json:"character"`
}

// --- LspService ---

type LspService struct{ t *transport.Client }

func (s *LspService) ep(id, endpoint string) string {
	return fmt.Sprintf("/v1/sandbox/%s/lsp/%s", id, endpoint)
}

func (s *LspService) Start(ctx context.Context, id string, p *LspServerParams) (*LspStatus, error) {
	var r LspStatus
	return &r, wrapErr("lsp start", s.t.JSON(ctx, http.MethodPost, s.ep(id, "start"), p, &r, nil))
}

func (s *LspService) Stop(ctx context.Context, id string, p *LspServerParams) (*LspStatus, error) {
	var r LspStatus
	return &r, wrapErr("lsp stop", s.t.JSON(ctx, http.MethodPost, s.ep(id, "stop"), p, &r, nil))
}

func (s *LspService) DidOpen(ctx context.Context, id string, p *LspDocumentParams) error {
	return wrapErr("lsp didOpen", s.t.JSON(ctx, http.MethodPost, s.ep(id, "did-open"), p, nil, nil))
}

func (s *LspService) DidClose(ctx context.Context, id string, p *LspDocumentParams) error {
	return wrapErr("lsp didClose", s.t.JSON(ctx, http.MethodPost, s.ep(id, "did-close"), p, nil, nil))
}

func (s *LspService) Completions(ctx context.Context, id string, p *LspPositionParams) (*LspResult, error) {
	var r LspResult
	return &r, wrapErr("lsp completions", s.t.JSON(ctx, http.MethodPost, s.ep(id, "completions"), p, &r, nil))
}

func (s *LspService) DocumentSymbols(ctx context.Context, id string, p *LspDocumentParams) (*LspResult, error) {
	var r LspResult
	return &r, wrapErr("lsp document symbols", s.t.JSON(ctx, http.MethodPost, s.ep(id, "document-symbols"), p, &r, nil))
}
