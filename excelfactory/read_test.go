package excelfactory

import (
	"testing"
)

const (
	filepath = "/Users/christianhovenbitzer/Desktop/rentabilität01-10.xlsx"
)

func TestRead(t *testing.T) {
	testFilterValue := "AIJI"

	file, err := ReadFile(filepath)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	for _, sh := range file.Sheets {

		// test FilterRowsByColumn
		vv, err := sh.FilterRowsByColumn("B", func(value string) bool { return value == testFilterValue })
		if err != nil {
			t.Errorf("got error %s", err)
		}
		for _, row := range vv {
			t.Log(row)
		}

		// test filterColumn
		v, err := sh.FilterColumn("B", func(value string) bool { return value == testFilterValue })
		if err != nil {
			t.Errorf("got error %s", err)
		}
		for _, cell := range v {
			if cell != testFilterValue {
				t.Errorf("got %s, expected %s", cell, testFilterValue)
			}
			t.Log(cell)
		}
	}
}
