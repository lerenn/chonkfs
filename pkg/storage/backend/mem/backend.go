package mem

import (
	"context"
	"strings"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

var _ backend.BackEnd = (*backEnd)(nil)

type backEnd struct {
	root *directory
}

func NewBackEnd() *backEnd {
	return &backEnd{
		root: newDirectory(),
	}
}

func (b *backEnd) CreateDirectory(_ context.Context, path string) error {
	path = strings.Trim(path, "/")
	return b.root.CreateDirectory(path)
}

func (b *backEnd) IsDirectory(_ context.Context, path string) error {
	path = strings.Trim(path, "/")
	return b.root.IsDirectory(path)
}

func (b *backEnd) CreateFile(_ context.Context, path string, chunkSize int) error {
	path = strings.Trim(path, "/")
	return b.root.CreateFile(path, chunkSize)
}

func (b *backEnd) IsFile(ctx context.Context, path string) error {
	path = strings.Trim(path, "/")
	return b.root.IsFile(path)
}

func (b *backEnd) ListFiles(ctx context.Context, path string) ([]string, error) {
	path = strings.Trim(path, "/")
	return b.root.ListFiles(path)
}
