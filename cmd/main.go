package main

import (
	"os"

	"github.com/saromanov/mystery/internal/mystery"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func put(c *cli.Context) error {
	key := c.Args().Get(0)
	value := c.Args().Get(1)
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	if err := mystery.Put(mystery.PutRequest{
		MasterPass: masterPass,
		Key:        key,
		Value:      value,
	}); err != nil {
		logrus.WithError(err).Fatalf("unable to store data")
	}
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
