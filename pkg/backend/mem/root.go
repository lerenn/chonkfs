package mem

import (
	"github.com/lerenn/chonkfs/pkg/backend"
)

var _ backend.Root = (*root)(nil)

type root struct {
	*directory
}

func New() *root {
	return &root{
		directory: newEmptyDirectory(),
	}
}
