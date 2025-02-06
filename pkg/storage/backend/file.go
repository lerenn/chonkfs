package backend

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
)

type File interface {
	GetInfo(ctx context.Context) (info.File, error)
}
