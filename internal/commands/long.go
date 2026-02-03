package commands

import (
	"fmt"
	"my-ls/internal/models"
	"os"
	"syscall"
)

func PrintLong(entries []models.Entry) {
	if len(entries) > 1 {
		fmt.Printf("total %d\n", totalBlocks(entries))
	}
	var maxLinks, maxSize int
	var maxUser, maxGroup int

	for _, e := range entries {
		st := stat(e)
		if st == nil {
			continue
		}

		links := int(st.Nlink)
		size := int(e.Info.Size())

		u := username(st.Uid)
		g := groupname(st.Gid)

		if digits(links) > maxLinks {
			maxLinks = digits(links)
		}
		if digits(size) > maxSize {
			maxSize = digits(size)
		}
		if len(u) > maxUser {
			maxUser = len(u)
		}
		if len(g) > maxGroup {
			maxGroup = len(g)
		}
	}

	for _, e := range entries {
		st := stat(e)
		if st == nil {
			continue
		}

		fmt.Printf(
			"%s %*d %-*s %-*s %*d %s %s\n",
			e.Info.Mode().String(),
			maxLinks, st.Nlink,
			maxUser, username(st.Uid),
			maxGroup, groupname(st.Gid),
			maxSize, e.Info.Size(),
			formatTime(e.Info.ModTime()),
			formatName(e),
		)
	}

}

func stat(e models.Entry) *syscall.Stat_t {
	if st, ok := e.Info.Sys().(*syscall.Stat_t); ok {
		return st
	}
	return nil
}

func totalBlocks(entries []models.Entry) int64 {
	var total int64

	for _, e := range entries {
		st := stat(e)
		if st == nil {
			continue
		}
		total += st.Blocks
	}

	return total
}

func formatName(e models.Entry) string {
	if e.Info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(e.Path); err == nil {
			return e.Name + " -> " + target
		}
	}
	return e.Name
}
