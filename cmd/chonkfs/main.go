package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/lerenn/chonkfs/pkg/backends/mem"
	"github.com/lerenn/chonkfs/pkg/chonkfs"
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

		// Create backend
		backend := mem.NewDirectory(
			mem.WithDirectoryLogger(logger))

		// Create chonkfs
		chFS := chonkfs.NewDirectory(backend,
			chonkfs.WithDirectoryLogger(logger),
			chonkfs.WithDirectoryChunkSize(chunkSize))

		// Create server
		server, err := fs.Mount(path, chFS, &fs.Options{
			Logger: logger,
			UID:    uint32(os.Getuid()),
			GID:    uint32(os.Getgid()),
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
	rootCmd.PersistentFlags().IntVarP(&chunkSize, "chunk-size", "s", chonkfs.DefaultChunkSize, "Set chunk size")

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s", err.Error())
		errCode = 1
	}

	// Exit with error code
	os.Exit(errCode)
}
