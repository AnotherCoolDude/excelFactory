package excelfactory

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

// New creates a new File
func New(path string) *File {
	f := excelize.NewFile()
	return &File{
		path: path,
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

		// fill sheet with data
		for rowIdx, row := range sh.data {
			for colIdx, cell := range row {
				coords, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
				if err != nil {
					return err
				}
				err = file.file.SetCellStr(sheetname, coords, cell)
				if err != nil {
					return err
				}
			}
		}
	}
	return file.file.SaveAs(file.path)
}

// SaveAs saves file at path
func (file *File) SaveAs(path string) error {
	file.path = path
	return file.Save()
}
