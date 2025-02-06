package backend

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
)

type Directory interface {
	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetInfo(ctx context.Context) (info.Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)

	CreateFile(ctx context.Context, name string, chunkSize int) (File, error)
	GetFile(ctx context.Context, name string) (File, error)
	ListFiles(ctx context.Context) (map[string]File, error)
}

type DirectoryInfo struct {
}
