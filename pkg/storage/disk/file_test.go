package disk

import (
	"context"
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
	Path string
}

func (suite *FileSuite) SetupTest() {
	path, err := os.MkdirTemp("", "chonkfs-test-*")
	suite.Require().NoError(err)
	suite.Path = path
	suite.Directory = NewDirectory(path)
}

func (suite *FileSuite) TearDownTest() {
	err := os.RemoveAll(suite.Path)
	suite.Require().NoError(err)
}

func (suite *FileSuite) TestWriteChunkData() {
	f, err := suite.Directory.CreateFile(context.Background(), "File-TestWriteChunkData.txt", info.File{
		ChunkSize: 8,
	})
	suite.Require().NoError(err)

	// Resize chunk count
	err = f.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Write a chunk
	buf := []byte("Hello")
	written, err := f.WriteChunk(context.Background(), 0, buf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), written)

	// Check the chunk is the good size
	stats, err := os.Stat(f.(*file).getChunckPath(0))
	suite.Require().NoError(err)
	suite.Require().Equal(int64(8), stats.Size())

	// Write outside the chunk
	buf = []byte("HelloWorld!")
	written, err = f.WriteChunk(context.Background(), 0, buf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(8, written)

	// Check the chunk is still the good size
	stats, err = os.Stat(f.(*file).getChunckPath(0))
	suite.Require().NoError(err)
	suite.Require().Equal(int64(8), stats.Size())
}
