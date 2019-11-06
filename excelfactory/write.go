package excelfactory

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

// New creates a new File
func New(path string) *File {
	f := excelize.NewFile()
	return &File{
		Path: path,
		file: f,
		Sheets: []Sheet{
			Sheet{
				data:          [][]string{},
				Name:          "Sheet1",
				HeaderColumns: map[string]string{},
			},
		},
	}
}

// AppendData appends data to sheet
func (sheet *Sheet) AppendData(data [][]string) {
	// append data
	sheet.data = append(sheet.data, data...)

	// evaluate new MaxCol
	for _, row := range data {
		if len(row) > sheet.MaxCol {
			sheet.MaxCol = len(row)
		}
	}

	// evaluate new MaxRow
	sheet.MaxRow = len(sheet.data)
}

// Save writes file to path
func (file *File) Save() error {
	sheetmap := file.file.GetSheetMap()

	for _, sh := range file.Sheets {
		sheetname := ""

		// does sheet exist
		for _, shname := range sheetmap {
			if shname == sh.Name {
				sheetname = shname
			}
		}

		// if sheetname is empty, create a new sheet (sheet doesnt exist)
		if sheetname == "" {
			file.file.NewSheet(sh.Name)
			sheetname = sh.Name
		}

		// define modifier to calculate coordinates
		rowModifier := 1
		colModifier := 1

		// check if sheet has header defined, write header row to file and adjust modifier
		if len(sh.HeaderColumns) != 0 {
			for col, hname := range sh.HeaderColumns {
				err := file.file.SetCellStr(sh.Name, fmt.Sprintf("%s%d", col, 1), hname)
				if err != nil {
					return err
				}
			}
			rowModifier++
		}

		// fill sheet with data
		for rowIdx, row := range sh.data {

			for colIdx, cell := range row {
				coords, err := excelize.CoordinatesToCellName(colIdx+colModifier, rowIdx+rowModifier)
				if err != nil {
					return err
				}
				// try to parse to float
				f, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					// if parsing fails, write as string
					err = file.file.SetCellStr(sheetname, coords, cell)
				} else {
					// if parsing succeeds, write as float
					err = file.file.SetCellFloat(sh.Name, coords, f, 2, 64)
				}
				if err != nil {
					return err
				}
			}
		}

		// format data
		for _, finfo := range sh.finfo {
			sID, err := file.file.NewStyle(finfo.format)
			if err != nil {
				return fmt.Errorf("could not parse format %s: %s", finfo.format, err)
			}
			splitArea := strings.Split(finfo.area, ":")
			err = file.file.SetCellStyle(sh.Name, splitArea[0], splitArea[1], sID)
			if err != nil {
				return fmt.Errorf("could not set columnstyle: %s", err)
			}
		}

	}
	return file.file.SaveAs(file.Path)
}

// SaveAs saves file at path
func (file *File) SaveAs(path string) error {
	file.Path = path
	return file.Save()
}
