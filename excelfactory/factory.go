package excelfactory

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

// File wraps excelize.File in a convenient struct
type File struct {
	path   string
	file   *excelize.File
	Sheets []Sheet
}

// Sheet wraps excelize.Sheet in a convenient struct
type Sheet struct {
	data          [][]string
	Name          string
	MaxCol        int
	MaxRow        int
	HeaderColumns map[string]string
}
