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
	/*var data mystery.Data
	for i := 1; i < c.Args().Len(); i++ {
		data += mystery.Data(c.Args().Get(i) + ";")
	}*/
	data, storeType := prepareData(c)
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
		Type:       storeType,
	}); err != nil {
		return fmt.Errorf("unable to store data: %v", err)
	}
	return nil
}

func prepareData(c cli.Args) (mystery.Data, string) {
	var data mystery.Data
	first := c.Args().Get(1)
	if strings.HasPrefix(first, "@") {
		return mystery.Data(), "file"
	}
	for i := 1; i < c.Args().Len(); i++ {
		data += mystery.Data(c.Args().Get(i) + ";")
	}
	return data, "store"
}

func get(c *cli.Context) error {
	conf, err := loadConfig("config.yml")
	if err != nil {
		log.WithError(err).Fatalf("unable to load config")
	}
	key := c.Args().Get(0)
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	pg, err := postgres.New(conf)
	if err != nil {
		log.WithError(err).Fatalf("unable to init backend")
	}
	value, err := mystery.Get(mystery.GetRequest{
		MasterPass: masterPass,
		Namespace:  key,
		Backend:    pg,
	})
	if err != nil {
		log.WithError(err).Fatalf("unable to get data")
	}

	fmt.Println(value)
	return nil
}

// update provides updating of the value
func update(c *cli.Context) error {
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
	var data mystery.Data
	for i := 1; i < c.Args().Len(); i++ {
		data += mystery.Data(c.Args().Get(i) + ";")
	}
	if err = mystery.Update(mystery.UpdateRequest{
		MasterPass: masterPass,
		Namespace:  key,
		Data:       data,
		Backend:    pg,
	}); err != nil {
		logrus.WithError(err).Fatalf("unable to update data")
	}
	return nil
}

func delete(c *cli.Context) error {
	conf, err := loadConfig("config.yml")
	if err != nil {
		logrus.WithError(err).Fatalf("unable to load config")
	}
	key := c.Args().Get(0)
	if key == "" {
		logrus.WithError(err).Fatalf("key is not defined")
	}
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	pg, err := postgres.New(conf)
	if err != nil {
		logrus.WithError(err).Fatalf("unable to init backend")
	}
	if err = mystery.Delete(mystery.DeleteRequest{
		MasterPass: masterPass,
		Namespace:  key,
		Backend:    pg,
	}); err != nil {
		logrus.WithError(err).Fatalf("unable to delete data")
	}

	log.Infof("data %s was deleted", key)
	return nil
}

func list(c *cli.Context) error {
	conf, err := loadConfig("config.yml")
	if err != nil {
		logrus.WithError(err).Fatalf("unable to load config")
	}
	masterPass := os.Getenv("MYSTERY_MASTER_PASS")
	pg, err := postgres.New(conf)
	if err != nil {
		logrus.WithError(err).Fatalf("unable to init backend")
	}
	rsp, err := mystery.List(mystery.ListRequest{
		MasterPass: masterPass,
		Backend:    pg,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("unable to list data")
	}

	for _, r := range rsp {
		fmt.Println("KEY: ", r.Namespace, r.CurrentVersion)
	}
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
			{
				Name:   "update",
				Usage:  "updating of the value by the key",
				Action: update,
			},
			{
				Name:   "delete",
				Usage:  "delete value by the key",
				Action: delete,
			},
			{
				Name:   "list",
				Usage:  "showing list of secrets",
				Action: list,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
