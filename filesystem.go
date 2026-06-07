package leap0

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// --- Filesystem Types ---

type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`
	Size        int64  `json:"size"`
	Permissions string `json:"permissions,omitempty"`
	ModTime     string `json:"mod_time,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Group       string `json:"group,omitempty"`
}

type DirListing struct {
	Files []FileInfo `json:"files"`
}

type TreeEntry struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Children []TreeEntry `json:"children,omitempty"`
}

type Tree struct {
	Tree []TreeEntry `json:"tree"`
}

type GrepMatch struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

type FileEdit struct {
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

type EditResult struct {
	Applied int    `json:"applied"`
	Content string `json:"content,omitempty"`
}

type MultiEditParams struct {
	Paths   []string `json:"paths"`
	Find    string   `json:"find"`
	Replace string   `json:"replace,omitempty"`
}

type MultiEditResult struct {
	ModifiedFiles []string `json:"modified_files"`
	TotalMatches  int      `json:"total_matches"`
}

type WriteFileEntry struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	Permissions string `json:"permissions,omitempty"`
}

type ReadFileParams struct {
	Offset int
	Limit  int
	Head   int
	Tail   int
}

// --- FilesystemService ---

type FilesystemService struct{ t *transport.Client }

func (s *FilesystemService) ep(id, endpoint string) string {
	return fmt.Sprintf("/v1/sandbox/%s/filesystem/%s", id, endpoint)
}

func (s *FilesystemService) Ls(ctx context.Context, id, path string, recursive bool, exclude []string) (*DirListing, error) {
	body := map[string]any{"path": path}
	if recursive {
		body["recursive"] = true
	}
	if len(exclude) > 0 {
		body["exclude"] = exclude
	}
	var r DirListing
	return &r, wrapErr("ls", s.t.JSON(ctx, http.MethodPost, s.ep(id, "ls"), body, &r, nil))
}

func (s *FilesystemService) Stat(ctx context.Context, id, path string) (*FileInfo, error) {
	var r FileInfo
	return &r, wrapErr("stat", s.t.JSON(ctx, http.MethodPost, s.ep(id, "stat"), map[string]any{"path": path}, &r, nil))
}

func (s *FilesystemService) Mkdir(ctx context.Context, id, path string, recursive bool) error {
	return wrapErr("mkdir", s.t.JSON(ctx, http.MethodPost, s.ep(id, "mkdir"), map[string]any{"path": path, "recursive": recursive}, nil, nil))
}

func (s *FilesystemService) WriteFile(ctx context.Context, id, path, content, permissions string) error {
	body := map[string]any{"path": path, "content": content}
	if permissions != "" {
		body["permissions"] = permissions
	}
	return wrapErr("write file", s.t.JSON(ctx, http.MethodPost, s.ep(id, "write-file"), body, nil, nil))
}

func (s *FilesystemService) WriteFiles(ctx context.Context, id string, files []WriteFileEntry) error {
	return wrapErr("write files", s.t.JSON(ctx, http.MethodPost, s.ep(id, "write-files"), map[string]any{"files": files}, nil, nil))
}

func (s *FilesystemService) ReadFile(ctx context.Context, id, path string, p *ReadFileParams) (string, error) {
	body := map[string]any{"path": path}
	if p != nil {
		if p.Offset > 0 {
			body["offset"] = p.Offset
		}
		if p.Limit > 0 {
			body["limit"] = p.Limit
		}
		if p.Head > 0 {
			body["head"] = p.Head
		}
		if p.Tail > 0 {
			body["tail"] = p.Tail
		}
	}
	r, err := s.t.Text(ctx, http.MethodPost, s.ep(id, "read-file"), body, &transport.Options{Headers: map[string]string{"Accept": "text/plain"}})
	return r, wrapErr("read file", err)
}

func (s *FilesystemService) ReadBytes(ctx context.Context, id, path string, p *ReadFileParams) ([]byte, error) {
	body := map[string]any{"path": path}
	if p != nil {
		if p.Offset > 0 {
			body["offset"] = p.Offset
		}
		if p.Limit > 0 {
			body["limit"] = p.Limit
		}
		if p.Head > 0 {
			body["head"] = p.Head
		}
		if p.Tail > 0 {
			body["tail"] = p.Tail
		}
	}
	r, err := s.t.Bytes(ctx, http.MethodPost, s.ep(id, "read-file"), body, &transport.Options{Headers: map[string]string{"Accept": "application/octet-stream"}})
	return r, wrapErr("read bytes", err)
}

func (s *FilesystemService) ReadFiles(ctx context.Context, id string, paths []string) (map[string]string, error) {
	var r map[string]string
	err := s.t.JSON(ctx, http.MethodPost, s.ep(id, "read-files"), map[string]any{"paths": paths}, &r, nil)
	return r, wrapErr("read files", err)
}

func (s *FilesystemService) Delete(ctx context.Context, id, path string, recursive bool) error {
	return wrapErr("delete", s.t.JSON(ctx, http.MethodPost, s.ep(id, "delete"), map[string]any{"path": path, "recursive": recursive}, nil, nil))
}

func (s *FilesystemService) SetPermissions(ctx context.Context, id, path, mode, owner, group string) error {
	body := map[string]any{"path": path}
	if mode != "" {
		body["mode"] = mode
	}
	if owner != "" {
		body["owner"] = owner
	}
	if group != "" {
		body["group"] = group
	}
	return wrapErr("set permissions", s.t.JSON(ctx, http.MethodPost, s.ep(id, "set-permissions"), body, nil, nil))
}

func (s *FilesystemService) Glob(ctx context.Context, id, path, pattern string, exclude []string) ([]string, error) {
	body := map[string]any{"path": path, "pattern": pattern}
	if len(exclude) > 0 {
		body["exclude"] = exclude
	}
	var r []string
	return r, wrapErr("glob", s.t.JSON(ctx, http.MethodPost, s.ep(id, "glob"), body, &r, nil))
}

func (s *FilesystemService) Grep(ctx context.Context, id, path, pattern string, include, exclude []string) ([]GrepMatch, error) {
	body := map[string]any{"path": path, "pattern": pattern}
	if len(include) > 0 {
		body["include"] = include
	}
	if len(exclude) > 0 {
		body["exclude"] = exclude
	}
	var r []GrepMatch
	return r, wrapErr("grep", s.t.JSON(ctx, http.MethodPost, s.ep(id, "grep"), body, &r, nil))
}

func (s *FilesystemService) EditFile(ctx context.Context, id, path string, edits []FileEdit) (*EditResult, error) {
	var r EditResult
	return &r, wrapErr("edit file", s.t.JSON(ctx, http.MethodPost, s.ep(id, "edit-file"), map[string]any{"path": path, "edits": edits}, &r, nil))
}

func (s *FilesystemService) EditFiles(ctx context.Context, id string, p *MultiEditParams) (*MultiEditResult, error) {
	var r MultiEditResult
	return &r, wrapErr("edit files", s.t.JSON(ctx, http.MethodPost, s.ep(id, "edit-files"), p, &r, nil))
}

func (s *FilesystemService) Move(ctx context.Context, id, src, dst string, overwrite bool) error {
	return wrapErr("move", s.t.JSON(ctx, http.MethodPost, s.ep(id, "move"), map[string]any{"src_path": src, "dst_path": dst, "overwrite": overwrite}, nil, nil))
}

func (s *FilesystemService) Copy(ctx context.Context, id, src, dst string, recursive, overwrite bool) error {
	return wrapErr("copy", s.t.JSON(ctx, http.MethodPost, s.ep(id, "copy"), map[string]any{"src_path": src, "dst_path": dst, "recursive": recursive, "overwrite": overwrite}, nil, nil))
}

func (s *FilesystemService) Exists(ctx context.Context, id, path string) (bool, error) {
	var r struct {
		Exists bool `json:"exists"`
	}
	err := s.t.JSON(ctx, http.MethodPost, s.ep(id, "exists"), map[string]any{"path": path}, &r, nil)
	return r.Exists, wrapErr("exists", err)
}

func (s *FilesystemService) Tree(ctx context.Context, id, path string, maxDepth int, exclude []string) (*Tree, error) {
	body := map[string]any{"path": path}
	if maxDepth > 0 {
		body["max_depth"] = maxDepth
	}
	if len(exclude) > 0 {
		body["exclude"] = exclude
	}
	var r Tree
	return &r, wrapErr("tree", s.t.JSON(ctx, http.MethodPost, s.ep(id, "tree"), body, &r, nil))
}
