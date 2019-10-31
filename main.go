package main

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/etats"
	"github.com/AnotherCoolDude/excelFactory/excelfactory"
	"strconv"
)

const (
	filepath = "/Users/christianhovenbitzer/Desktop/rentabilität01-10.xlsx"
)

func main() {
	f, err := excelfactory.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	e, err := etats.ReadFile("../etats/etatdirector.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, d := range e.Directors {

	}

	rows, err := f.Sheets[0].FilterRowsByColumn("B", func(value string) bool {
		return value == "ABCH"
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	revenue := 0.0
	for _, row := range rows {
		for col, cell := range row {
			if col != 11 {
				continue
			}
			amount, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			revenue += amount
		}
	}
	// todo put revenue into ob
}
