package commands

import (
	"bufio"
	"os"
)

var out = bufio.NewWriter(os.Stdout)

func Flush() {
	out.Flush()
}
