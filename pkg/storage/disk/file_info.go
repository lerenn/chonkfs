package disk

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"

	"github.com/lerenn/chonkfs/pkg/storage"
)

func writeFileInfo(fi storage.FileInfo, filePath string) error {
	fileInfoPath := path.Join(filePath, fileInfoName)

	// Transform to json
	data, err := json.Marshal(fi)
	if err != nil {
		return err
	}

	// Write the file
	if err := os.WriteFile(fileInfoPath, data, fs.FileMode(0755)); err != nil {
		return err
	}

	return nil
}

func readFileInfo(filePath string) (storage.FileInfo, error) {
	fileInfoPath := path.Join(filePath, fileInfoName)

	// Read the file
	data, err := os.ReadFile(fileInfoPath)
	if err != nil {
		return storage.FileInfo{}, err
	}

	// Unmarshal json
	var fi storage.FileInfo
	if err := json.Unmarshal(data, &fi); err != nil {
		return storage.FileInfo{}, err
	}

	return fi, nil
}
