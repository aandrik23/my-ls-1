package commands

import (
	"fmt"
	"my-ls/internal/models"
	"unicode/utf8"
)

func PrintShort(entries []models.Entry) {
	for _, e := range entries {
		fmt.Fprintln(out, e.Name)
	}
}

func PrintEntry(entry models.Entry, flags models.Flags) {
	if flags.Long {
		PrintLong([]models.Entry{entry})
		return
	}
	fmt.Fprintln(out, entry.Name)
}

func PrintDirHeader(name string) {
	fmt.Fprintf(out, "%s:\n", name)
}

func PrintColumns(entries []models.Entry) {
	if len(entries) == 0 {
		return
	}

	sep := 2 // ls uses 2 spaces between columns
	termWidth := terminalWidth()

	// 1) max display width (runes, not bytes)
	maxLen := 0
	for _, e := range entries {
		if l := utf8.RuneCountInString(e.Name); l > maxLen {
			maxLen = l
		}
	}

	// If even one column doesn't fit, fall back to one-per-line
	if maxLen > termWidth {
		for _, e := range entries {
			fmt.Fprintln(out, e.Name)
		}
		return
	}

	// If everything fits on one row, print one row (ls behavior)
	totalWidth := 0
	for i, e := range entries {
		totalWidth += utf8.RuneCountInString(e.Name)
		if i != len(entries)-1 {
			totalWidth += sep
		}
	}

	if totalWidth <= termWidth {
		for i, e := range entries {
			if i > 0 {
				fmt.Fprint(out, "  ")
			}
			fmt.Fprint(out, e.Name)
		}
		fmt.Fprintln(out)
		return
	}

	// 2) how many columns fit
	cols := (termWidth + sep) / (maxLen + sep)
	if cols < 1 {
		cols = 1
	}

	// 3) how many rows
	rows := (len(entries) + cols - 1) / cols

	// 4) print column-major
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := c*rows + r
			if i >= len(entries) {
				continue
			}

			name := entries[i].Name
			if c == cols-1 || i+rows >= len(entries) {
				fmt.Fprint(out, name)
			} else {
				fmt.Fprintf(out, "%-*s  ", maxLen, name)
			}
		}
		fmt.Fprintln(out)
	}
}
