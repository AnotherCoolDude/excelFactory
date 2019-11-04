package excelfactory

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strings"
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
	finfo         []formattingInfo
	Name          string
	MaxCol        int
	MaxRow        int
	HeaderColumns map[string]string
}

type formattingInfo struct {
	area   string
	format string
}

// styles
var (
	FormatEuro = `{"number_format": 219}`
)

// FormatColumns formats the given area with the provided format
func (sheet *Sheet) FormatColumns(columns string, format string) {
	splitColumns := strings.Split(columns, ":")
	area := ""
	hasHeader := len(sheet.HeaderColumns) > 0
	startingRow := 1
	if hasHeader {
		startingRow = 2
	}
	if len(splitColumns) > 1 {
		area = fmt.Sprintf("%s%d:%s%d", splitColumns[0], startingRow, splitColumns[1], sheet.MaxRow)
	} else {
		area = fmt.Sprintf("%s%d:%s%d", columns, startingRow, columns, sheet.MaxRow)
	}
	sheet.finfo = append(sheet.finfo, formattingInfo{area, format})
}
