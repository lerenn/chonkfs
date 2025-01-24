package chonker

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	Directory Directory
	suite.Suite
}

func (suite *FileSuite) SetupTest() {
	d, err := NewDirectory(context.Background(), mem.NewDirectory())
	suite.Require().NoError(err)
	suite.Directory = d
}

func (suite *FileSuite) TestReadWrite() {
	f, err := suite.Directory.CreateFile(context.Background(), "File-TestReadWrite.txt", 4)
	suite.Require().NoError(err)

	for i := 0; i < 10; i++ {
		// Write a chunk
		buf := []byte("Hello, world!")
		written, err := f.Write(context.Background(), buf, i, WriteOptions{})
		suite.Require().NoError(err)
		suite.Require().Equal(len(buf), written)

		// Read the chunk
		readBuf := make([]byte, len(buf))
		readBuf, err = f.Read(context.Background(), readBuf, i)
		suite.Require().NoError(err)
		suite.Require().Equal([]byte("Hello, world!"), readBuf)
		suite.Require().Equal(len(buf), len(readBuf))
	}
}

func (suite *FileSuite) TestReadOnChunkSize() {
	f, err := suite.Directory.CreateFile(context.Background(), "File-TestReadWrite.txt", 4)
	suite.Require().NoError(err)

	// Write a chunk
	buf := []byte("1234")
	written, err := f.Write(context.Background(), buf, 0, WriteOptions{})
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), written)

	// Read the chunk
	readBuf := make([]byte, len(buf))
	readBuf, err = f.Read(context.Background(), readBuf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal([]byte("1234"), readBuf)
	suite.Require().Equal(len(buf), len(readBuf))
}

func (suite *FileSuite) TestReadWithBiggerBufferThanData() {
	f, err := suite.Directory.CreateFile(context.Background(), "File-TestReadWrite.txt", 4)
	suite.Require().NoError(err)

	// Write a chunk
	buf := []byte("1234")
	written, err := f.Write(context.Background(), buf, 0, WriteOptions{})
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), written)

	// Read the chunk
	readBuf := make([]byte, 10)
	readBuf, err = f.Read(context.Background(), readBuf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal([]byte("1234"), readBuf[:len(buf)])
	suite.Require().Equal(len(buf), len(readBuf))
}
