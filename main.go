package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	filepath = "/Users/christianhovenbitzer/Desktop/rentabilität01-10.xlsx"
)

func main() {
	file, err := readFile(filepath)
	if err != nil {
		fmt.Println(err)
		panic(0)
	}
	for _, sh := range file.Sheets {

		// test FilterRowsByColumn
		vv, err := sh.FilterRowsByColumn("B", func(value string) bool { return value == "AIJI" })
		if err != nil {
			fmt.Println(err)
			panic(0)
		}
		for _, row := range vv {
			fmt.Println(row)
		}

		// test filterColumn
		v, err := sh.FilterColumn("B", func(value string) bool { return value == "AIJI" })
		if err != nil {
			fmt.Println(err)
			panic(0)
		}
		for _, cell := range v {
			fmt.Println(cell)
		}
	}
}

// File wraps excelize.File in a convenient struct
type File struct {
	file   *excelize.File
	Sheets []Sheet
}

// Sheet wraps excelize.Sheet in a convenient struct
type Sheet struct {
	rows       [][]string
	Name       string
	MaxCol     int
	MaxRow     int
	HeaderCols map[string]string
}

func readFile(path string) (File, error) {
	file := File{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		return file, err
	}
	file.file = f

	names := make([]string, f.SheetCount)
	for _, sheetname := range f.GetSheetMap() {
		names = append(names, sheetname)
		sheet := Sheet{
			HeaderCols: map[string]string{},
		}
		sheet.Name = sheetname

		rows, err := f.GetRows(sheetname)
		if err != nil {
			return file, err
		}
		sheet.rows = rows

		for row, rowCells := range rows {
			sheet.MaxRow = len(rows)
			for col, cell := range rowCells {
				currentCoords, err := excelize.CoordinatesToCellName(col+1, row+1)
				if err != nil {
					return file, err
				}
				// set columns
				if row == 0 {
					sheet.HeaderCols[currentCoords] = cell
				}
				// set MaxCol if it's longer than before
				if sheet.MaxCol < len(rowCells) {
					sheet.MaxCol = len(rowCells)
				}

			}
		}

		//add sheet to file
		file.Sheets = append(file.Sheets, sheet)
	}
	return file, nil
}

// Filterfunc provides the current cell values und returns a bool
type Filterfunc func(cell string) bool

// FilterColumn filters all values in column using filter
func (sheet *Sheet) FilterColumn(column string, filter Filterfunc) ([]string, error) {
	values := []string{}
	vv, err := sheet.FilterRowsByColumn(column, filter)
	if err != nil {
		return values, err
	}
	for _, v := range vv {
		for col, cell := range v {
			colnr, _ := excelize.ColumnNameToNumber(column)
			if col != colnr-1 {
				continue
			}
			values = append(values, cell)
		}
	}
	return values, nil
}

// FilterRowsByColumn takes the column (e.g. "A") and a Filterfunc and returns all rows, that have a filtered columnvalue in column
func (sheet *Sheet) FilterRowsByColumn(column string, filter Filterfunc) ([][]string, error) {
	values := [][]string{}
	colnr, err := excelize.ColumnNameToNumber(column)
	if err != nil {
		return values, err
	}

	for row, cells := range sheet.rows {
		//skip header col and check if column exists
		if row == 0 {
			if colnr > len(cells) {
				return values, fmt.Errorf("column %s out of bounds", column)
			}
			continue
		}

		for col, cell := range cells {
			// find column to filter
			if col+1 != colnr {
				continue
			}
			// filter cell
			if !filter(cell) {
				continue
			}
			values = append(values, cells)
		}
	}
	return values, nil
}
