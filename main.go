package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// Create a new CLI app
	app := &cli.App{
		Name:  "probe",
		Usage: "Analyze a file and display insights.",
		Flags: []cli.Flag{}, // No global flags
		Action: func(c *cli.Context) error {
			// Check if the required argument (filename) is provided
			if c.NArg() < 1 {
				cli.ShowAppHelpAndExit(c, 1) // Display help and exit with error code
			}

			filename := c.Args().Get(0) // First argument is the filename
			fmt.Printf("Processing file: %s\n", filename)

			// Here you can add your file analysis logic
			return nil
		},
	}

	// Run the app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
