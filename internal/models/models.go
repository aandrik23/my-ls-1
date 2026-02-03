package models

import "os"

type Flags struct {
	Long      bool // -l
	All       bool // -a
	Recursive bool // -R
	Dir       bool // -d
	One       bool // -1
}

type Entry struct {
	Name string
	Path string
	Info os.FileInfo
}
