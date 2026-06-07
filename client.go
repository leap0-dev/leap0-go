package leap0

import "github.com/leap0dev/leap0-go/internal/transport"

// Client is the top-level entry point.
type Client struct {
	t *transport.Client

	Sandboxes       *SandboxesService
	Snapshots       *SnapshotsService
	Templates       *TemplatesService
	Filesystem      *FilesystemService
	Git             *GitService
	Process         *ProcessService
	Pty             *PtyService
	Lsp             *LspService
	SSH             *SSHService
	CodeInterpreter *CodeInterpreterService
	Desktop         *DesktopService
}

// NewClient creates a Leap0 client.
func NewClient(opts ...Option) (*Client, error) {
	_, tc, err := resolveConfig(opts)
	if err != nil {
		return nil, err
	}
	t := transport.New(tc)
	return &Client{
		t:               t,
		Sandboxes:       &SandboxesService{t: t},
		Snapshots:       &SnapshotsService{t: t},
		Templates:       &TemplatesService{t: t},
		Filesystem:      &FilesystemService{t: t},
		Git:             &GitService{t: t},
		Process:         &ProcessService{t: t},
		Pty:             &PtyService{t: t},
		Lsp:             &LspService{t: t},
		SSH:             &SSHService{t: t},
		CodeInterpreter: &CodeInterpreterService{t: t},
		Desktop:         &DesktopService{t: t},
	}, nil
}

// Close releases resources.
func (c *Client) Close() {}
