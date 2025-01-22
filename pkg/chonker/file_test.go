package chonker

import (
	"context"
	"testing"

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
	suite.Directory = NewDirectory()
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
		err = f.Read(context.Background(), readBuf, i)
		suite.Require().NoError(err)
		suite.Require().Equal([]byte("Hello, world!"), readBuf)
		suite.Require().Equal(len(buf), len(readBuf))
	}
}
