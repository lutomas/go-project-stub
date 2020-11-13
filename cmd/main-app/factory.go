package main

import (
	"fmt"

	"github.com/lutomas/go-project-stub/api/server"
	"github.com/lutomas/go-project-stub/internal/store/sql"
	"github.com/lutomas/go-project-stub/pkg/config"
	"go.uber.org/zap"
)

// Factory method to build store from configuration.
func MakeStoreFromConfig(cfg *config.Database, log *zap.Logger) (*sql.SQL, error) {
	opts := &sql.Options{
		Logger: log,
	}
	opts.DatabaseType = sql.DatabaseTypeSqlite3
	opts.URI = cfg.Sqlite3DB
	opts.MaxOpenConnections = cfg.MaxOpenConnections

	if cfg.PostgresHost != "" {
		pgHost := cfg.PostgresHost
		pgUser := cfg.PostgresUser
		pgPass := cfg.PostgresPassword
		pgDBName := cfg.PostgresDB

		opts.MaxIdleConnections = cfg.PostgresMaxIdleConnections

		opts.DatabaseType = sql.DatabaseTypePostgres
		opts.URI = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s", pgHost, pgUser, pgDBName, pgPass, cfg.TZ)
	}

	s, err := sql.NewSQL(opts)
	if err != nil {
		return nil, fmt.Errorf("DB setup failed: %s", err)
	}

	log.Info("database configured", zap.String("type", opts.DatabaseType))

	return s, nil
}

func MakeMainAppService(cfg *config.MainAppServer, log *zap.Logger) (*server.Server, error) {
	opts := &server.Options{
		HttpHost:   cfg.HttpHost,
		HttpPort:   cfg.HttpPort,
		Logger:     log,
		EnableCors: cfg.EnableCors,
	}

	return server.New(opts)
}
