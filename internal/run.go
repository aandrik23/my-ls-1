package internal

import (
	"fmt"
	"my-ls/internal/commands"
	"my-ls/internal/models"
	"my-ls/logger"
)

func Run(paths []string, flags models.Flags) {
	// If -d: treat everything as a file (DispatchPath already does this correctly)
	if flags.Dir {
		for _, p := range paths {
			if err := DispatchPath(p, flags); err != nil {
				logger.PrintError(logger.ACCESS, p, err)
			}
		}
	}

	files, dirs, errs := splitOperands(paths)
	for _, e := range errs {
		logger.PrintError(logger.ACCESS, e.Path, e.Err)
	}

	sortEntries(files, flags)
	sortEntries(dirs, flags)

	// Print files first
	for _, e := range files {
		commands.PrintEntry(e, flags)
	}

	// Then list directories
	needHeader := flags.Recursive || len(dirs) > 1 || len(files) > 0

	for i, d := range dirs {
		if needHeader {
			if i > 0 {
				fmt.Println()
			}
			commands.PrintDirHeader(d.Name)
		}

		if err := listDirectory(d.Path, d.Name, flags); err != nil {
			logger.PrintError(logger.ACCESS, d.Name, err)
		}
	}
	defer commands.Flush()
}
