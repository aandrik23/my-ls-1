package commands

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func digits(n int) int {
	if n == 0 {
		return 1
	}
	c := 0
	for n > 0 {
		n /= 10
		c++
	}
	return c
}

func username(uid uint32) string {
	u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return strconv.FormatUint(uint64(uid), 10)
	}
	return u.Username
}

func groupname(gid uint32) string {
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))
	if err != nil {
		return strconv.FormatUint(uint64(gid), 10)
	}
	return g.Name
}

func formatTime(t time.Time) string {
	now := time.Now()

	// 6 months ≈ 180 days
	sixMonths := time.Hour * 24 * 180

	if t.After(now) || now.Sub(t) > sixMonths {
		// old or future file → show year
		return t.Format("Jan _2  2006")
	}

	// recent file → show time
	return t.Format("Jan _2 15:04")
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func terminalWidth() int {
	ws := &winsize{}

	// stdout fd = 1
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(os.Stdout.Fd()),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if err != 0 || ws.Col == 0 {
		return 80 // fallback, just like ls
	}

	return int(ws.Col)
}
