package internal

import (
	"my-ls/internal/models"
	"os"
)

func DispatchPath(path string, flags models.Flags) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	entry := models.Entry{
		Name: path,
		Path: path,
		Info: info,
	}

	// -d means: treat directory like a file
	if info.IsDir() && !flags.Dir {
		return listDirectory(path, entry.Name, flags)
	}

	return listFile(entry, flags)
}
