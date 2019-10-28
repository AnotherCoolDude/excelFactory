package excelfactory

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	path := "/Users/christianhovenbitzer/Test.xlsx"

	nf := New(path)

	dd1 := [][]string{
		[]string{"Test 1", " Test 2", "Test 3", "Test 4"},
		[]string{"Value 1", "Value 2", "Value 3", "Value 4"},
	}

	dd2 := [][]string{
		[]string{"Value 1", " Value 2", "Value 3", "Value 4"},
		[]string{"Value 3", "Value 4", "Value 5", "Value 6"},
	}

	nf.Sheets[0].AppendData(dd1)
	nf.Sheets[0].AppendData(dd2)

	err := nf.Save()
	if err != nil {
		t.Errorf("could not save file, got error: %s", err)
		return
	}

	f, err := ReadFile(path)
	if err != nil {
		t.Errorf("could not open file, got error: %s", err)
		return
	}

	if len(f.Sheets[0].data) == 0 {
		t.Error("no data was saved to file")
	}

	for _, row := range f.Sheets[0].data {
		for _, cell := range row {
			t.Logf("%s\t", cell)
		}
		t.Log()
	}

	err = os.Remove(path)
	if err != nil {
		t.Errorf("could not delete file, got error: %s", err)
	}

}
