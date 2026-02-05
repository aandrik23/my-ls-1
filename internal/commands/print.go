package commands

import (
	"fmt"
	"my-ls/internal/models"
	"unicode/utf8"
)

func PrintShort(entries []models.Entry) {
	for _, e := range entries {
		fmt.Fprintln(out, formatColoredName(e))
	}
}

func PrintEntry(entry models.Entry, flags models.Flags) {
	if flags.Long {
		PrintLong([]models.Entry{entry})
		return
	}
	fmt.Fprintln(out, formatColoredName(entry))
}

func PrintDirHeader(name string) {
	fmt.Fprintf(out, "%s:\n", name)
}

func PrintColumns(entries []models.Entry) {
	if len(entries) == 0 {
		return
	}

	const colSep = 2 // ls uses 2 spaces between columns
	termWidth := terminalWidth()

	// Precompute display widths
	widths := make([]int, len(entries))
	maxLen := 0
	for i, e := range entries {
		w := utf8.RuneCountInString(e.Name)
		widths[i] = w
		if w > maxLen {
			maxLen = w
		}
	}

	// If even one column doesn't fit, fall back to one-per-line
	if maxLen > termWidth {
		for _, e := range entries {
			fmt.Fprintln(out, formatColoredName(e))
		}
		return
	}

	// If everything fits on one row, print one row (ls behavior)
	totalWidth := 0
	for i, w := range widths {
		totalWidth += w
		if i != len(widths)-1 {
			totalWidth += colSep
		}
	}

	if totalWidth <= termWidth {
		for i, e := range entries {
			if i > 0 {
				fmt.Fprint(out, "  ")
			}
			fmt.Fprint(out, formatColoredName(e))
		}
		fmt.Fprintln(out)
		return
	}

	// How many columns fit
	cols := (termWidth + colSep) / (maxLen + colSep)
	if cols < 1 {
		cols = 1
	}

	// How many rows
	rows := (len(entries) + cols - 1) / cols

	// Print column-major (ls style)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := c*rows + r
			if i >= len(entries) {
				continue
			}

			if c == cols-1 || i+rows >= len(entries) {
				fmt.Fprint(out, formatColoredName(entries[i]))
			} else {
				fmt.Fprintf(out, "%-*s  ", maxLen, formatColoredName(entries[i]))
			}
		}
		fmt.Fprintln(out)
	}
}
