package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/backends"
	"github.com/stretchr/testify/suite"
)

type FileSuite struct {
	Directory backends.Directory
	suite.Suite
}

func (suite *FileSuite) TestReadWrite() {
	f, err := suite.Directory.CreateFile(context.Background(), "File-TestReadWrite.txt", 4)
	suite.Require().NoError(err)

	for i := 0; i < 10; i++ {
		// Write a chunk
		buf := []byte("Hello, world!")
		written, err := f.Write(context.Background(), buf, i, backends.WriteOptions{})
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
