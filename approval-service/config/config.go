package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/server/common/log"
	"go.uber.org/zap"
)

type Cfg struct {
	RestPort       int
	MySQLConfig    MySQLConfig
	TemporalConfig TemporalConfig
}

type MySQLConfig struct {
	Username     string
	Password     string
	Protocol     string
	Address      string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

type TemporalConfig struct {
	HostPort string
}

func Load() *Cfg {
	v := viper.New()
	v.AutomaticEnv()

	cfg := &Cfg{
		RestPort: v.GetInt("REST_PORT"),
		TemporalConfig: TemporalConfig{
			HostPort: v.GetString("TEMPORAL_HOST_PORT"),
		},
		MySQLConfig: MySQLConfig{
			Username:     v.GetString("MYSQL_USERNAME"),
			Password:     v.GetString("MYSQL_PASSWORD"),
			Protocol:     v.GetString("MYSQL_PROTOCOL"),
			Address:      v.GetString("MYSQL_ADDRESS"),
			Database:     v.GetString("MYSQL_DATABASE"),
			MaxOpenConns: v.GetInt("MYSQL_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("MYSQL_MAX_IDLE_CONNS"),
		},
	}

	return cfg
}

func GetMySQLDSN(cfg MySQLConfig) string {
	return fmt.Sprintf(
		"%s:%s@%s(%s)/%s?parseTime=true&multiStatements=true",
		cfg.Username,
		cfg.Password,
		cfg.Protocol,
		cfg.Address,
		cfg.Database,
	)
}

func GetTemporalClientOptions(cfg TemporalConfig, logger *zap.Logger) client.Options {
	return client.Options{
		HostPort: cfg.HostPort,
		Logger:   log.NewSdkLogger(log.NewZapLogger(logger)),
	}
}
