package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func put(c *cli.Context) error {
	return nil
}

func main() {
	app := &cli.App{
		Name:  "mystery",
		Usage: "Starting of the app",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:   "put",
				Usage:  "putting of key-value pair",
				Action: put,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
