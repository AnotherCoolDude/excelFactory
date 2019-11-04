package main

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/etats"
	"github.com/AnotherCoolDude/excelFactory/excelfactory"
	"github.com/urfave/cli"
	"os"
	"path"
	"strconv"
)

const (
	readexcelpath  = "rent01-09.xlsx"
	writeexcelpath = "result.xlsx"
	jsonpath       = "etats/etatdirector.json"
)

func main() {
	app := cli.NewApp()
	app.Name = "orderbook"
	app.Usage = "takes a rentabilität from proad, and sorts it to fit into orderbooks"

	app.Action = func(c *cli.Context) error {
		proadpath := c.Args().First()

		// check if path is a xlsx file
		if path.Ext(proadpath) != ".xlsx" {
			return fmt.Errorf("provided path needs to be a xlsx file\n %s", proadpath)
		}
		// run task
		runTask(proadpath, jsonpath)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("error running program: %s\n", err)
		os.Exit(1)
	}

}

func runTask(rentpath, etatpath string) {
	// get excel file
	f, err := excelfactory.ReadFile(rentpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get json file
	e, err := etats.ReadFile(etatpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create new file
	newFile := excelfactory.New(path.Join(path.Dir(rentpath), "etatsummary.xlsx"))

	// loop over directors
	for _, d := range e.Directors {
		// loop over campaigns
		for _, camp := range d.Campaigns {
			// get all rows for camp
			rows, err := f.Sheets[0].FilterRowsByColumn("B", func(value string) bool {
				return value == camp
			})

			if err != nil {
				fmt.Printf("error filtering excel file: %s\n", err)
				os.Exit(1)
			}

			// check if rows are available
			if len(rows) == 0 {
				fmt.Printf("Campaign %s: Could not find any entries in file %s\n", camp, path.Base(readexcelpath))
				continue
			}

			// add up income of rows
			income := 0.0
			for _, row := range rows {
				val, err := strconv.ParseFloat(row[11], 64)
				if err != nil {
					fmt.Printf("error parsing float %.2f: %s\n", val, err)
					os.Exit(1)
				}
				income += val
			}
			newFile.Sheets[0].AppendData([][]string{{d.Name, camp, strconv.FormatFloat(income, 'f', 2, 64)}})
		}
	}

	//set style
	newFile.Sheets[0].FormatColumns("C", excelfactory.FormatEuro)

	// check if all camps in excel file where read
	unreadCamp := ""
	rows, err := f.Sheets[0].FilterColumn("B", func(value string) bool {
		for _, d := range e.Directors {
			for _, c := range d.Campaigns {
				unreadCamp = value
				return c == value
			}
		}
		return false
	})
	if err != nil {
		fmt.Printf("error filtering excel file: %s\n", err)
		os.Exit(1)
	}
	if len(rows) == 0 {
		fmt.Printf("campaign %s was not used\n", unreadCamp)
	}

	// save file
	err = newFile.Save()
	if err != nil {
		fmt.Printf("could not save file to path %s: %s\n", writeexcelpath, err)
	}
}
