package mem

import (
	"github.com/lerenn/chonkfs/pkg/backends"
)

var _ backends.Root = (*root)(nil)

type root struct {
	*directory
}

func New() *root {
	return &root{
		directory: newEmptyDirectory(),
	}
}
