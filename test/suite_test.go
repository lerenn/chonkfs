package test

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/chonker"
	fuse1 "github.com/lerenn/chonkfs/pkg/fuse"
	"github.com/lerenn/chonkfs/pkg/storage/layer"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

const (
	testDir = "/tmp/chonkfs"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	mountPoints []*fuse.Server
	suite.Suite
}

func (suite *Suite) TearDownSuite() {
	// Unmount all mount points
	for _, srv := range suite.mountPoints {
		_ = srv.Unmount()
	}

	// Clean up the suite playground
	suiteName := strings.Split(suite.T().Name(), "/")[0]
	err := os.RemoveAll(testDir + "/" + suiteName)
	suite.Require().NoError(err)
}

func (suite *Suite) createChonkFS(
	backend chonker.Directory,
	chunkSize int,
) (path string, server *fuse.Server) {
	// Create a directory corresponding to the test name
	path = testDir + "/" + suite.T().Name()
	err := os.MkdirAll(path, os.ModePerm)
	suite.Require().NoError(err)

	// Create a chonkfs
	chFS := fuse1.NewDirectory(backend,
		fuse1.WithDirectoryChunkSize(chunkSize))

	// Mount the ChonkFS
	server, err = fs.Mount(path, chFS, &fs.Options{
		UID: uint32(os.Getuid()),
		GID: uint32(os.Getgid()),
	})
	suite.Require().NoError(err)

	// Append the mount point to the list
	suite.mountPoints = append(suite.mountPoints, server)

	return
}

func (suite *Suite) TestWriteOnlyThenReadOnly() {
	// Mount chunkfs
	c, err := chonker.NewDirectory(context.Background(), layer.NewDirectory(mem.NewDirectory(), nil))
	suite.Require().NoError(err)
	path, srv := suite.createChonkFS(c, 4096)

	// --- WRITE FILE

	// Create file
	f, err := os.OpenFile(path+"/hello.txt", os.O_WRONLY|os.O_CREATE, 0755)
	suite.Require().NoError(err)

	// Write to file
	buf := []byte("Hello, World!")
	n, err := f.Write(buf)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Close file
	err = f.Close()
	suite.Require().NoError(err)

	// --- READ FILE

	// Open file
	f, err = os.OpenFile(path+"/hello.txt", os.O_RDONLY, 0755)
	suite.Require().NoError(err)

	// Read from file
	readBuf := make([]byte, len(buf))
	n, err = f.Read(readBuf)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Close file
	err = f.Close()
	suite.Require().NoError(err)

	// -- COMPARE

	// Compare buffers
	suite.Require().Equal(string(buf), string(readBuf))

	// Unmount chunkfs
	err = srv.Unmount()
	suite.Require().NoError(err)
}

func (suite *Suite) TestReadWriteMode() {
	// Mount chunkfs
	c, err := chonker.NewDirectory(context.Background(), layer.NewDirectory(mem.NewDirectory(), nil))
	suite.Require().NoError(err)
	path, srv := suite.createChonkFS(c, 4096)

	// Create file
	f, err := os.OpenFile(path+"/hello.txt", os.O_RDWR|os.O_CREATE, 0755)
	suite.Require().NoError(err)

	// Write to file
	buf := []byte("Hello, World!")
	n, err := f.Write(buf)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Return to the beginning of the file
	_, err = f.Seek(0, 0)
	suite.Require().NoError(err)

	// Read from file
	readBuf := make([]byte, len(buf))
	n, err = f.Read(readBuf)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Compare buffers
	suite.Require().Equal(string(buf), string(readBuf))

	// Close file
	err = f.Close()
	suite.Require().NoError(err)

	// Unmount chunkfs
	err = srv.Unmount()
	suite.Require().NoError(err)
}

func (suite *Suite) TestRandomReadWrite() {
	// Mount chunkfs
	c, err := chonker.NewDirectory(context.Background(), layer.NewDirectory(mem.NewDirectory(), nil))
	suite.Require().NoError(err)
	path, srv := suite.createChonkFS(c, 4096)

	// Create file
	f, err := os.OpenFile(path+"/hello.txt", os.O_RDWR|os.O_CREATE, 0755)
	suite.Require().NoError(err)

	// Loop through random read/write operations
	buf := []byte("hello")
	for i := 0; i < 100; i++ {
		pos := rand.Intn(100)

		// Write buffer at random position
		w, err := f.WriteAt(buf, int64(pos))
		suite.Require().NoError(err)
		suite.Require().Equal(len(buf), w)

		// Read buffer at same position
		readBuf := make([]byte, len(buf))
		r, err := f.ReadAt(readBuf, int64(pos))
		suite.Require().NoError(err)
		suite.Require().Equal(len(buf), r)

		// Compare buffers
		suite.Require().Equal(string(buf), string(readBuf))
	}

	// Close file
	err = f.Close()
	suite.Require().NoError(err)

	// Unmount chunkfs
	err = srv.Unmount()
	suite.Require().NoError(err)
}
