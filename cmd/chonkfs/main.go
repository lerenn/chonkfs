package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/lerenn/chonkfs/pkg/chonker"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/wrapper"
	"github.com/spf13/cobra"
)

var (
	path      string
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
		if path == "" {
			return fmt.Errorf("path is required")
		}

		// Create a default logger if logging is activated
		var logger *log.Logger
		if debug {
			logger = log.Default()
		} else {
			logger = log.New(os.Stdout, "", 0)
		}

		// Create chonker
		c, err := chonker.NewDirectory(
			cmd.Context(),
			mem.NewDirectory(nil),
			chonker.WithDirectoryLogger(logger))
		if err != nil {
			return err
		}

		// Create wrapper for FUSE
		w := wrapper.NewDirectory(c,
			wrapper.WithDirectoryLogger(logger),
			wrapper.WithDirectoryChunkSize(chunkSize))

		// Create FUSE server
		to := time.Duration(1)
		server, err := fs.Mount(path, w, &fs.Options{
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
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Set mount path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().IntVarP(&chunkSize, "chunk-size", "s", wrapper.DefaultChunkSize, "Set chunk size")

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s", err.Error())
		errCode = 1
	}

	// Exit with error code
	os.Exit(errCode)
}
