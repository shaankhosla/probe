package main

import (
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"probe/cmd"
)

func main() {
	app := &cli.App{
		Name:  "probe",
		Usage: "Interactive SQL query tool for file analysis.",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("error: you must provide a filename as an argument")
			}
			filename := c.Args().Get(0)

			return cmd.RunProbe(filename)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
