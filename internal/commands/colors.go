package commands

import "my-ls/internal/models"

const (
	colorReset = "\033[0m"
	colorBlue  = "\033[34m"
)

func formatColoredName(e models.Entry) string {
	if e.Info.IsDir() {
		return colorBlue + e.Name + colorReset
	}
	return e.Name
}
