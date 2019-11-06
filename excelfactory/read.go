package excelfactory

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// ReadFile reads an excel file at path and prepares it
func ReadFile(path string, header bool) (File, error) {
	file := File{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		return file, err
	}
	file.file = f
	file.Path = path

	names := make([]string, f.SheetCount)
	for _, sheetname := range f.GetSheetMap() {
		names = append(names, sheetname)
		sheet := Sheet{
			HeaderColumns: map[string]string{},
		}
		sheet.Name = sheetname

		rows, err := f.GetRows(sheetname)
		if err != nil {
			return file, err
		}

		if header {
			for col, cell := range rows[0] {
				colname, err := excelize.ColumnNumberToName(col + 1)
				if err != nil {
					return file, err
				}
				sheet.HeaderColumns[colname] = cell
			}
			sheet.data = rows[1:]
		} else {
			sheet.data = rows
		}

		for _, rowCells := range sheet.data {
			sheet.MaxRow = len(sheet.data)

			// set MaxCol if it's longer than before
			if sheet.MaxCol < len(rowCells) {
				sheet.MaxCol = len(rowCells)
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

	// check if column exists
	if len(sheet.data[0]) < colnr {
		return values, fmt.Errorf("column %s out of bounds", column)
	}

	for _, row := range sheet.data {
		for col, cell := range row {
			// find column to filter
			if col+1 != colnr {
				continue
			}
			// filter cell
			if !filter(cell) {
				continue
			}
			values = append(values, row)
		}
	}
	return values, nil
}
