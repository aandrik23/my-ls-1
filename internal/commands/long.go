package commands

import (
	"fmt"
	"my-ls/internal/models"
	"os"
	"syscall"
)

var userCache = make(map[uint32]string)
var groupCache = make(map[uint32]string)

func PrintLong(entries []models.Entry) {
	if len(entries) > 1 {
		fmt.Fprintf(out, "total %d\n", totalBlocks(entries))
	}
	var maxLinks, maxSize int
	var maxUser, maxGroup int

	for _, e := range entries {
		st, ok := statOK(e)
		if !ok {
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
		st, ok := statOK(e)
		if !ok {
			continue
		}

		fmt.Fprintf(out,
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

func statOK(e models.Entry) (*syscall.Stat_t, bool) {
	return e.Stat, e.Stat != nil
}

func totalBlocks(entries []models.Entry) int64 {
	var total int64

	for _, e := range entries {
		st, ok := statOK(e)
		if !ok {
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
