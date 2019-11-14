package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "excel factory"
	app.Usage = "a library / cli tool to quickly read from and create excel files"
	app.Version = "0.1.0"

	app.Commands = []*cli.Command{
		orderbookCommand,
		transferCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("error running program: %s\n", err)
		os.Exit(1)
	}

}
