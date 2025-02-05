package mem

import (
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

type file struct {
	chunkSize int
}

func newFile(chunkSize int) (*file, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("%w: %d", backend.ErrInvalidChunkSize, chunkSize)
	}

	return &file{
		chunkSize: chunkSize,
	}, nil
}
