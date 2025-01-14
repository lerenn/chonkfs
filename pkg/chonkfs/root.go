package chonkfs

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/lerenn/chonkfs/pkg/backends"
)

// Capabilities that the root struct should implements
var _ fs.NodeOnAdder = (*Root)(nil)

type Root struct {
	directory
}

func NewRoot(backend backends.Root) *Root {
	return &Root{
		directory: directory{
			backendDirectory: backend,
		},
	}
}

func (r *Root) OnAdd(ctx context.Context) {
	// TODO
}
