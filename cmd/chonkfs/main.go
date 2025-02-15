package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/lerenn/chonkfs/pkg/chonker"
	"github.com/lerenn/chonkfs/pkg/fuse"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/disk"
	"github.com/lerenn/chonkfs/pkg/storage/layer"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/spf13/cobra"
)

var (
	diskPath  string
	mntPath   string
	debug     bool
	chunkSize int
)

var rootCmd = &cobra.Command{
	Use:     "chonkfs",
	Version: "0.1.0",
	Short:   "chonkfs - a CLI to manage ChonkFS",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Check if path is set
		if mntPath == "" {
			return fmt.Errorf("mount path is required")
		}

		// Create a default logger if logging is activated
		var logger *log.Logger
		if debug {
			logger = log.Default()
		} else {
			logger = log.New(io.Discard, "", 0)
		}

		// Create backend
		var err error
		var be storage.Directory
		switch {
		case diskPath != "":
			be, err = layer.NewDirectory(mem.NewDirectory(), disk.NewDirectory(diskPath))
			if err != nil {
				return err
			}
		default:
			be = mem.NewDirectory()
		}

		// Create chonker
		c, err := chonker.NewDirectory(cmd.Context(), be, chonker.WithDirectoryLogger(logger))
		if err != nil {
			return err
		}

		// Create wrapper for FUSE
		w := fuse.NewDirectory(c,
			fuse.WithDirectoryLogger(logger),
			fuse.WithDirectoryChunkSize(chunkSize))

		// Create FUSE server
		to := time.Duration(1)
		server, err := fs.Mount(mntPath, w, &fs.Options{
			Logger:       logger,
			UID:          uint32(os.Getuid()),
			GID:          uint32(os.Getgid()),
			EntryTimeout: &to,
			AttrTimeout:  &to,
		})
		if err != nil {
			return err
		}

		// Wait for server to finish
		server.Wait()
		return nil
	},
}

func main() {
	var errCode int

	// Set flags
	rootCmd.PersistentFlags().StringVarP(&mntPath, "mnt", "m", "", "Set mount path")
	rootCmd.PersistentFlags().StringVarP(&diskPath, "disk", "d", "", "Set disk path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().IntVarP(&chunkSize, "chunk-size", "s", fuse.DefaultChunkSize, "Set chunk size")

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s", err.Error())
		errCode = 1
	}

	// Exit with error code
	os.Exit(errCode)
}
