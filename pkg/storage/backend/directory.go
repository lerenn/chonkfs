package backend

import "context"

type Directory interface {
	CreateDirectory(ctx context.Context, name string) error
	IsDirectory(ctx context.Context, name string) error

	CreateFile(ctx context.Context, name string, chunkSize int) error
	IsFile(ctx context.Context, name string) error
	ListFiles(ctx context.Context) ([]string, error)
}
