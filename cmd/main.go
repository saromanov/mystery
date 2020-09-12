package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/backend/postgres"
	"github.com/saromanov/mystery/internal/mystery"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func put(c *cli.Context) error {
	if err := putInner(c); err != nil {
		log.WithError(err).Fatalf("unable to execute put command")
	}
	log.Infof("data was stored")
	return nil
}

func putInner(c *cli.Context) error {
	conf, err := loadConfig("config.yml")
	if err != nil {
		return fmt.Errorf("unable to load config: %v", err)
	}
	key := c.Args().Get(0)
	data := mystery.Data{}
	for i := 1; i < c.Args().Len(); i++ {
		value := strings.Split(c.Args().Get(i), "=")
		if len(value) <= 1 {
			return fmt.Errorf("data should be in format key=value")
		}
		data[value[0]] = value[1]
	}
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	pg, err := postgres.New(conf)
	if err != nil {
		return fmt.Errorf("unable to init backend: %v", err)
	}
	if err := mystery.Put(mystery.PutRequest{
		MasterPass: masterPass,
		Namespace:  key,
		Data:       data,
		Backend:    pg,
	}); err != nil {
		return fmt.Errorf("unable to store data: %v", err)
	}
	return nil
}

func get(c *cli.Context) error {
	conf, err := loadConfig("config.yml")
	if err != nil {
		logrus.WithError(err).Fatalf("unable to load config")
	}
	key := c.Args().Get(0)
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	pg, err := postgres.New(conf)
	if err != nil {
		logrus.WithError(err).Fatalf("unable to init backend")
	}
	value, err := mystery.Get(mystery.GetRequest{
		MasterPass: masterPass,
		Namespace:  key,
		Backend:    pg,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("unable to get data")
	}

	fmt.Println(value)
	return nil
}

// loadConfig provides loading of configuration
func loadConfig(path string) (*config.Config, error) {
	return config.Load(path)
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
			{
				Name:   "get",
				Usage:  "getting value by the key",
				Action: get,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
