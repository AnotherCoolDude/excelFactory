package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/AnotherCoolDude/excelFactory/etats"
	"github.com/AnotherCoolDude/excelFactory/excelfactory"
	"github.com/urfave/cli"
)

// define name for resulting file
const (
	resultFilename = "result.xlsx"
)

// define flag variables and command
var (
	header   bool
	rentpath string
	jsonpath string
	// define command to use in main.go
	orderbookCommand = &cli.Command{
		Name:    "orderbook",
		Aliases: []string{"ob"},
		Usage:   "takes a rentabilität xlsx from proad and a etat json file and sorts it to fit into orderbooks",
		Action:  orderbookAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proad, p",
				Usage:       "path to proad rentabilität xlsx file",
				Destination: &rentpath,
			},
			&cli.StringFlag{
				Name:        "etat, e",
				Usage:       "path to etat json file",
				Value:       "etats/etatdirector.json",
				Destination: &jsonpath,
			},
			&cli.BoolFlag{
				Name:        "header",
				Usage:       "indicate wether the provided xlsx file has a header column, true by default",
				Destination: &header,
			},
		},
	}
)

// wrap task in action to use in orderbook command
func orderbookAction(c *cli.Context) error {
	// check if rentpath is a xlsx file
	if path.Ext(rentpath) != ".xlsx" {
		return fmt.Errorf("provided path needs to be a xlsx file\n %s", rentpath)
	}

	// check if jsonpath is a json file
	if path.Ext(jsonpath) != ".json" {
		return fmt.Errorf("provided path needs to be a json file\n %s", jsonpath)
	}

	// run task
	obTask(rentpath, jsonpath, header)
	return nil
}

// the actual task to perform
func obTask(rentpath, jsonpath string, header bool) {
	// get excel file
	f, err := excelfactory.ReadFile(rentpath, header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get json file
	e, err := etats.ReadFile(jsonpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create new file
	newFile := excelfactory.New(path.Join(path.Dir(rentpath), "etatsummary.xlsx"))

	// set header row
	newFile.Sheets[0].HeaderColumns = map[string]string{"A": "Manager", "B": "Kampagne", "C": "Erlös", "D": "Summe Erlöse"}

	// loop over directors
	for _, d := range e.Directors {
		// loop over campaigns
		for i, camp := range d.Campaigns {
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
				fmt.Printf("campaign %s: Could not find any entries in file %s\n", camp, path.Base(rentpath))
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
			// if end of camps from d has been reached
			if len(d.Campaigns)-1 == i {
				rows, _ := newFile.Sheets[0].FilterRowsByColumn("A", func(value string) bool { return value == d.Name })
				summary := 0.0
				for _, row := range rows {
					val, err := strconv.ParseFloat(row[2], 64)
					if err != nil {
						fmt.Println(err)
					}
					summary += val
				}
				summary += income

				// add a summary of d in new column
				newFile.Sheets[0].AppendData([][]string{{d.Name, camp, strconv.FormatFloat(income, 'f', 2, 64), strconv.FormatFloat(summary, 'f', 2, 64)}})
			} else {
				// simply add a row with name, campaign and income
				newFile.Sheets[0].AppendData([][]string{{d.Name, camp, strconv.FormatFloat(income, 'f', 2, 64)}})
			}
		}
	}

	//set style
	newFile.Sheets[0].FormatColumns("C", excelfactory.FormatEuro)
	newFile.Sheets[0].FormatColumns("D", excelfactory.FormatEuro)

	// check if all camps in excel file where read
	rows, err := f.Sheets[0].FilterRowsByColumn("B", func(value string) bool {
		existing := false
		for _, d := range e.Directors {
			for _, c := range d.Campaigns {
				if c == value {
					existing = true
				}
			}
		}
		// filter only non existing values
		return !existing
	})
	if err != nil {
		fmt.Printf("error filtering excel file: %s\n", err)
		os.Exit(1)
	}

	for _, row := range rows {
		fmt.Printf("campaign %s from file %s was not used\n", row[1], path.Base(rentpath))
		newFile.Sheets[0].AppendData([][]string{{"Unbekannt", row[1], row[11]}})
	}

	// save file
	err = newFile.Save()
	if err != nil {
		fmt.Printf("could not save file to path %s: %s\n", newFile.Path, err)
	}
}

// helper func to eliminate duplicate strings in slice
func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
