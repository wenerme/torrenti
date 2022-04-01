package serve

//go:generate gomodifytags -file=conf.go -w -all -add-tags yaml -transform snakecase --skip-unexported -add-options yaml=omitempty

import (
	"time"

	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type GORMConf struct {
	DisableForeignKeyConstraintWhenMigrating bool       `yaml:"disable_foreign_key_constraint_when_migrating,omitempty"`
	Log                                      SQLLogConf `envPrefix:"LOG_" yaml:"log,omitempty"`
}
type DatabaseConf struct {
	Type         string   `env:"TYPE" yaml:"type,omitempty"`
	Driver       string   `env:"DRIVER" yaml:"driver,omitempty"`
	Database     string   `env:"DATABASE" yaml:"database,omitempty"`
	Username     string   `env:"USERNAME" yaml:"username,omitempty"`
	Password     string   `env:"PASSWORD" yaml:"password,omitempty"`
	Host         string   `env:"HOST" yaml:"host,omitempty"`
	Port         int      `env:"PORT" yaml:"port,omitempty"`
	Schema       string   `env:"SCHEMA" yaml:"schema,omitempty"`
	CreateSchema bool     `env:"CREATE_SCHEMA" yaml:"create_schema,omitempty"`
	DSN          string   `env:"DSN" yaml:"dsn,omitempty"`
	GORM         GORMConf `envPrefix:"GORM_" yaml:"gorm,omitempty"`

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
type GRPClientCConf struct {
	Addr string `env:"ADDR" yaml:"addr,omitempty"`
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
