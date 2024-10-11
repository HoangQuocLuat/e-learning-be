package main

import (
	"context"
	"e-learning/config"
	"e-learning/src/cronjob"
	"e-learning/src/database"
	face_config "e-learning/src/face-config"
	kafka_config "e-learning/src/kafka"
	"e-learning/src/server"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	configPrefix string
	configSource string
)

const dataDir = "/home/ad/Documents/e-learning/e-learning-be/train"

func main() {
	app := cli.NewApp()
	app.Name = "E-Learning microservice"
	app.Usage = "E-Learning microservice"
	app.Copyright = "Copyright © 2024 HoangLuat. All Rights Reserved."
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "configPrefix",
			Aliases:     []string{"confPrefix"},
			Usage:       "prefix for config",
			Value:       "learning",
			Destination: &configPrefix,
		},
		&cli.StringFlag{
			Name:        "configSource",
			Aliases:     []string{"confSource"},
			Value:       "../config/.env",
			Usage:       "set path to environment file",
			Destination: &configSource,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "Start the e-learning server",
			Action: Serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr-graph",
					Aliases: []string{"address-graph"},
					Value:   "0.0.0.0:8989",
					Usage:   "address for serve graph",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		return config.LoadFromEnv(configPrefix, configSource)
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endSignal := make(chan os.Signal, 1)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func(ctx context.Context, errChan chan error) {
		err := app.RunContext(ctx, os.Args)
		errChan <- err
	}(ctx, errChan)

	select {
	case sign := <-endSignal:
		log.Println("shutting down. reason:", sign)
		return
	case err := <-errChan:
		if err == nil {
			return
		}
		log.Println("encountered error:", err)
		return
	}
}

func Serve(c *cli.Context) error {
	ctx := c.Context
	err := database.ConnectDatabse(ctx)
	if err != nil {
		panic(err)
	}
	kafka_config.InitKafkaProducer()
	face_config.InitRecognizer(dataDir)
	go func() { cronjob.NotifyWithTimeBySchedules() }()
	go func() { cronjob.ComputeTuition()}()

	return server.ServeGraph(c.Context, c.String("addr-graph"))
}
