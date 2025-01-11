package main

import (
	"fmt"
	"os"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/lerenn/chonkfs/pkg/backend/mem"
	"github.com/lerenn/chonkfs/pkg/chonkfs"
	"github.com/spf13/cobra"
)

var (
	path string
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

		// Set debug mode to true
		chonkfs.SetDebug(true)

		// Create server
		// TODO: make backend configurable
		server, err := fs.Mount(path, chonkfs.NewRoot(mem.New()), &fs.Options{})
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

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s", err.Error())
		errCode = 1
	}

	// Exit with error code
	os.Exit(errCode)
}
