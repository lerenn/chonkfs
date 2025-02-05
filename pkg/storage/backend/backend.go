package backend

import "context"

type BackEnd interface {
	CreateDirectory(ctx context.Context, path string) error
	IsDirectory(ctx context.Context, path string) error

	CreateFile(ctx context.Context, path string, chunkSize int) error
	IsFile(ctx context.Context, path string) error
}
