package main

//go:generate gomodifytags -file=conf.go -w -all -add-tags yaml -transform snakecase --skip-unexported -add-options yaml=omitempty

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"go.uber.org/multierr"
)

type Config struct {
	util.DirConf `yaml:",inline"`
	PluginDir    string       `env:"PLUGIN_DIR" yaml:"plugin_dir,omitempty"`
	DB           DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
	Log          LogConf      `envPrefix:"LOG_" yaml:"log,omitempty"`
	HTTP         HTTPConf     `envPrefix:"HTTP_" yaml:"http,omitempty"`
	Debug        DebugConf    `envPrefix:"DEBUG_" yaml:"debug,omitempty"`
	GRPC         GRPCConf     `envPrefix:"GRPC_" yaml:"grpc,omitempty"`
	Scrape       ScrapeConf   `envPrefix:"SCRAPE_" yaml:"scrape,omitempty"`

	Torrent TorrentConf `envPrefix:"TORRENT_" yaml:"torrent,omitempty"`
	Sub     SubConf     `envPrefix:"SUB_" yaml:"sub,omitempty"`
}

func (conf *Config) defaults() {
	if conf.PluginDir == "" {
		self := os.Args[0]
		conf.PluginDir = filepath.Join(filepath.Dir(self), "plugins")

		if strings.Contains(self, "go-build") {
			conf.PluginDir = "bin/plugins"
		}
	}

	if conf.DB.Type != "" {
		m := map[string]interface{}{}
		err := mapstructure.Decode(conf.DB, &m)
		for k, v := range m {
			switch {
			case v == "":
			case v == 0:
			default:
				continue
			}
			delete(m, k)
		}
		err = multierr.Combine(
			err,
			mapstructure.Decode(m, &conf.Torrent.DB),
			mapstructure.Decode(m, &conf.Sub.DB),
			mapstructure.Decode(m, &conf.Scrape.Store.DB),
		)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	if conf.Torrent.DB.Type == "" {
		defaultTo(&conf.Sub.DB, DatabaseConf{
			Type:     "sqlite",
			Database: "$DATA_DIR/torrenti.sqlite",
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
	if conf.Sub.DB.Type == "" {
		defaultTo(&conf.Sub.DB, DatabaseConf{
			Type:     "sqlite",
			Database: "$DATA_DIR/subi.sqlite",
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
	if conf.Scrape.Store.DB.Type == "" {
		defaultTo(&conf.Scrape.Store.DB, DatabaseConf{
			Type:     "sqlite",
			Database: "$CACHE_DIR/scraper-store.db",
			Log: SQLLogConf{
				IgnoreNotFound: true,
			},
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
}

func defaultTo(a interface{}, def interface{}) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(def, &m)
	err = multierr.Combine(err, mapstructure.Decode(a, &m))
	for k, v := range m {
		del := false
		switch vv := v.(type) {
		case string:
			del = len(vv) == 0
		case int:
			del = vv == 0
		case map[string]interface{}:
			del = len(vv) == 0
		case []interface{}:
			del = len(vv) == 0
		}

		if del {
			delete(m, k)
		}
	}
	err = multierr.Combine(err, mapstructure.WeakDecode(m, a))
	if err != nil {
		log.Fatal().Err(err).Msg("apply default")
	}
}

type DatabaseConf struct {
	Type         string     `env:"TYPE" yaml:"type,omitempty"`
	Driver       string     `env:"DRIVER" yaml:"driver,omitempty"`
	Database     string     `env:"DATABASE" yaml:"database,omitempty"`
	Username     string     `env:"USERNAME" yaml:"username,omitempty"`
	Password     string     `env:"PASSWORD" yaml:"password,omitempty"`
	Host         string     `env:"HOST" yaml:"host,omitempty"`
	Port         int        `env:"PORT" yaml:"port,omitempty"`
	Schema       string     `env:"SCHEMA" yaml:"schema,omitempty"`
	CreateSchema bool       `env:"CREATE_SCHEMA" yaml:"create_schema,omitempty"`
	DSN          string     `env:"DSN" yaml:"dsn,omitempty"`
	Log          SQLLogConf `envPrefix:"LOG_" yaml:"log,omitempty"`

	DriverOptions DatabaseDriverOptions `envPrefix:"DRIVER_" yaml:"driver_options,omitempty"`
	Attributes    map[string]string     `envPrefix:"ATTR_" yaml:"attributes,omitempty"` // ConnectionAttributes
}

type DatabaseDriverOptions struct {
	MaxIdleConnections int            `env:"MAX_IDLE_CONNS" yaml:"max_idle_connections,omitempty"`
	MaxOpenConnections int            `env:"MAX_OPEN_CONNS" yaml:"max_open_connections,omitempty"`
	ConnMaxIdleTime    time.Duration  `env:"MAX_IDLE_TIME" yaml:"conn_max_idle_time,omitempty"`
	ConnMaxLifetime    *time.Duration `env:"MAX_LIVE_TIME" yaml:"conn_max_lifetime,omitempty"`
}

type SQLLogConf struct {
	SlowThreshold  time.Duration `env:"SLOW_THRESHOLD" yaml:"slow_threshold,omitempty"`
	IgnoreNotFound bool          `env:"IGNORE_NOT_FOUND" yaml:"ignore_not_found,omitempty"`
	Debug          bool          `env:"DEBUG" yaml:"debug,omitempty"`
}

type LogConf struct {
	Level string `env:"LEVEL" envDefault:"info" yaml:"level,omitempty"`
}

type GRPCConf struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool            `env:"ENABLED" envDefault:"true" yaml:"enabled,omitempty"`
	Gateway         GRPCGatewayConf `envPrefix:"GATEWAY_" yaml:"gateway,omitempty"`
}
type GRPCGatewayConf struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool   `env:"ENABLED" envDefault:"true" yaml:"enabled,omitempty"`
	Prefix          string `env:"PREFIX" yaml:"prefix,omitempty"`
}

type HTTPConf struct {
	util.ListenConf `yaml:",inline"`
}
type DebugConf struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool `env:"ENABLED" envDefault:"true" yaml:"enabled,omitempty"`
}
type TorrentConf struct {
	DB DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
}
type SubConf struct {
	DB DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
}

type ScrapeStoreConf struct {
	DB DatabaseConf `envPrefix:"STORE_DB_" yaml:"db,omitempty"`
}
type ScrapeConf struct {
	Debug ScrapeDebugConf `envPrefix:"DEBUG_" yaml:"debug,omitempty"`
	Store ScrapeStoreConf `envPrefix:"STORE_" yaml:"store,omitempty"`
}
type ScrapeDebugConf struct {
	Addr string `env:"ADDR" envDefault:"127.0.0.1:7676" yaml:"addr,omitempty"`
}
