package leap0

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Git Types ---

type GitResult struct {
	Output string `json:"output"`
}

type GitCommitResult struct {
	Output   string `json:"output"`
	CommitID string `json:"commit_id"`
}

type CloneParams struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Branch   string `json:"branch,omitempty"`
	CommitID string `json:"commit_id,omitempty"`
	Depth    int    `json:"depth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LogParams struct {
	Path           string `json:"path"`
	MaxCount       int    `json:"max_count,omitempty"`
	StartTimestamp int64  `json:"start_timestamp,omitempty"`
	EndTimestamp   int64  `json:"end_timestamp,omitempty"`
}

type BranchListParams struct {
	Path        string `json:"path"`
	BranchType  string `json:"branch_type,omitempty"`
	Contains    string `json:"contains,omitempty"`
	NotContains string `json:"not_contains,omitempty"`
}

type CreateBranchParams struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Checkout   bool   `json:"checkout,omitempty"`
	BaseBranch string `json:"base_branch,omitempty"`
}

type CheckoutParams struct {
	Path        string `json:"path"`
	Branch      string `json:"branch"`
	Create      bool   `json:"create,omitempty"`
	SetUpstream string `json:"set_upstream,omitempty"`
}

type DeleteBranchParams struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Force bool   `json:"force,omitempty"`
}

type CommitParams struct {
	Path       string `json:"path"`
	Message    string `json:"message"`
	Author     string `json:"author,omitempty"`
	Email      string `json:"email,omitempty"`
	AllowEmpty bool   `json:"allow_empty,omitempty"`
}

type PushParams struct {
	Path        string `json:"path"`
	Remote      string `json:"remote,omitempty"`
	Branch      string `json:"branch,omitempty"`
	SetUpstream bool   `json:"set_upstream,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
}

type PullParams struct {
	Path        string `json:"path"`
	Remote      string `json:"remote,omitempty"`
	Branch      string `json:"branch,omitempty"`
	Rebase      bool   `json:"rebase,omitempty"`
	SetUpstream bool   `json:"set_upstream,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
}

// --- GitService ---

type GitService struct{ t *transport.Client }

func (s *GitService) ep(id, endpoint string) string {
	return fmt.Sprintf("/v1/sandbox/%s/git/%s", id, endpoint)
}

func (s *GitService) Clone(ctx context.Context, id string, p *CloneParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git clone", s.t.JSON(ctx, http.MethodPost, s.ep(id, "clone"), p, &r, nil))
}

func (s *GitService) Status(ctx context.Context, id, path string) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git status", s.t.JSON(ctx, http.MethodPost, s.ep(id, "status"), map[string]any{"path": path}, &r, nil))
}

func (s *GitService) Branches(ctx context.Context, id string, p *BranchListParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git branches", s.t.JSON(ctx, http.MethodPost, s.ep(id, "branches"), p, &r, nil))
}

func (s *GitService) DiffUnstaged(ctx context.Context, id, path string, contextLines int) (*GitResult, error) {
	body := map[string]any{"path": path}
	if contextLines > 0 {
		body["context_lines"] = contextLines
	}
	var r GitResult
	return &r, wrapErr("git diff unstaged", s.t.JSON(ctx, http.MethodPost, s.ep(id, "diff-unstaged"), body, &r, nil))
}

func (s *GitService) DiffStaged(ctx context.Context, id, path string, contextLines int) (*GitResult, error) {
	body := map[string]any{"path": path}
	if contextLines > 0 {
		body["context_lines"] = contextLines
	}
	var r GitResult
	return &r, wrapErr("git diff staged", s.t.JSON(ctx, http.MethodPost, s.ep(id, "diff-staged"), body, &r, nil))
}

func (s *GitService) Diff(ctx context.Context, id, path, target string, contextLines int) (*GitResult, error) {
	body := map[string]any{"path": path, "target": target}
	if contextLines > 0 {
		body["context_lines"] = contextLines
	}
	var r GitResult
	return &r, wrapErr("git diff", s.t.JSON(ctx, http.MethodPost, s.ep(id, "diff"), body, &r, nil))
}

func (s *GitService) Reset(ctx context.Context, id, path string) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git reset", s.t.JSON(ctx, http.MethodPost, s.ep(id, "reset"), map[string]any{"path": path}, &r, nil))
}

func (s *GitService) Log(ctx context.Context, id string, p *LogParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git log", s.t.JSON(ctx, http.MethodPost, s.ep(id, "log"), p, &r, nil))
}

func (s *GitService) Show(ctx context.Context, id, path, revision string) (*GitResult, error) {
	body := map[string]any{"path": path}
	if revision != "" {
		body["revision"] = revision
	}
	var r GitResult
	return &r, wrapErr("git show", s.t.JSON(ctx, http.MethodPost, s.ep(id, "show"), body, &r, nil))
}

func (s *GitService) CreateBranch(ctx context.Context, id string, p *CreateBranchParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git create branch", s.t.JSON(ctx, http.MethodPost, s.ep(id, "create-branch"), p, &r, nil))
}

func (s *GitService) Checkout(ctx context.Context, id string, p *CheckoutParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git checkout", s.t.JSON(ctx, http.MethodPost, s.ep(id, "checkout-branch"), p, &r, nil))
}

func (s *GitService) DeleteBranch(ctx context.Context, id string, p *DeleteBranchParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git delete branch", s.t.JSON(ctx, http.MethodPost, s.ep(id, "delete-branch"), p, &r, nil))
}

func (s *GitService) Add(ctx context.Context, id, path string, files []string) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git add", s.t.JSON(ctx, http.MethodPost, s.ep(id, "add"), map[string]any{"path": path, "files": files}, &r, nil))
}

func (s *GitService) Commit(ctx context.Context, id string, p *CommitParams) (*GitCommitResult, error) {
	var r GitCommitResult
	return &r, wrapErr("git commit", s.t.JSON(ctx, http.MethodPost, s.ep(id, "commit"), p, &r, nil))
}

func (s *GitService) Push(ctx context.Context, id string, p *PushParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git push", s.t.JSON(ctx, http.MethodPost, s.ep(id, "push"), p, &r, nil))
}

func (s *GitService) Pull(ctx context.Context, id string, p *PullParams) (*GitResult, error) {
	var r GitResult
	return &r, wrapErr("git pull", s.t.JSON(ctx, http.MethodPost, s.ep(id, "pull"), p, &r, nil))
}
