package commands

import (
	"fmt"
	"my-ls/internal/models"
)

func PrintShort(entries []models.Entry) {
	for _, e := range entries {
		fmt.Println(e.Name)
	}
}

func PrintEntry(entry models.Entry, flags models.Flags) {
	if flags.Long {
		PrintLong([]models.Entry{entry})
		return
	}
	fmt.Println(entry.Name)
}

func PrintDirHeader(name string) {
	fmt.Printf("%s:\n", name)
}

const termWidth = 80

func PrintColumns(entries []models.Entry) {
	if len(entries) == 0 {
		return
	}

	// 1) find max name width
	maxLen := 0
	for _, e := range entries {
		if l := len(e.Name); l > maxLen {
			maxLen = l
		}
	}

	colWidth := maxLen + 1 // padding like ls
	width := terminalWidth()

	// 2) how many columns fit
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}
	if cols == 1 {
		for _, e := range entries {
			fmt.Println(e.Name)
		}
		return
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
				fmt.Print(name)
			} else {
				fmt.Printf("%-*s", colWidth, name)
			}
		}
		fmt.Println()
	}
}
