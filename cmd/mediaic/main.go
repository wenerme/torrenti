package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/wenerme/torrenti/pkg/serve"
)

const Name = "mediaic"

func main() {
	app := &cli.App{
		Name:  Name,
		Usage: "media indexer client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file",
				Value:   "",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "",
			},
		},
		Before: setup,
		Commands: cli.Commands{
			{
				Name: "stat",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "index",
				Subcommands: cli.Commands{
					{
						Name: "add",
						Action: func(context *cli.Context) error {
							return nil
						},
					},
				},
			},
		},
	}
	log.Err(app.Run(os.Args)).Send()
}

type Config struct {
	Log    serve.LogConf        `envPrefix:"LOG_" yaml:"log,omitempty"`
	Client serve.GRPClientCConf `envPrefix:"CLIENT_" yaml:"client,omitempty"`
}

func setup(ctx *cli.Context) (err error) {
	return
}

var _conf = &Config{
	Client: serve.GRPClientCConf{Addr: ":18080"},
}
