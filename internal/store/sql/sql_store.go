package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lutomas/go-project-stub/internal/store"
	"github.com/lutomas/go-project-stub/pkg/zap_logger"
	"github.com/lutomas/go-project-stub/types"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	// importing postgres driver
	_ "gorm.io/driver/postgres"
)

const DatabaseTypeSqlite3 = "sqlite3"
const DatabaseTypePostgres = "postgres"

type Options struct {
	DatabaseType       string // sqlite3 / postgres
	URI                string // path or conn string
	MaxOpenConnections int    // Max open connection
	MaxIdleConnections int    // Max idle connection
	// This flag, if it is `TRUE` -  will prevent db schema migration.
	SkipAutoMigrate  bool
	EnableSqlLogging bool

	Logger *zap.Logger
}

type SQL struct {
	db           *gorm.DB
	sqlDB        *sql.DB
	databaseType string
	logger       *zap.Logger
}

// To verify if SQL type matches store.Store interface.
var _ store.Store = &SQL{}

func NewSQL(opts *Options) (_ *SQL, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	if opts.Logger == nil {
		opts.Logger = zap_logger.GetInstance()
	}

	db, sqlDB, err := connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if db.Dialector.Name() == "sqlite" {
		db.Exec("PRAGMA foreign_keys = ON")
	}

	defer func() {
		if err != nil {
			// If error happens - close DB connections.
			sqlDB.Close()
		} else {
			// Enable SQL logging after DB schema migration
			if opts.EnableSqlLogging {
				fmt.Println("+++++++++ SQL-logging enabled")
				// db.LogMode(res)
				db.Logger = db.Logger.LogMode(gormLogger.Info)
			} else {
				db.Logger = db.Logger.LogMode(gormLogger.Silent)
			}
		}
	}()

	s := &SQL{
		db:     db,
		sqlDB:  sqlDB,
		logger: opts.Logger,
	}

	if !opts.SkipAutoMigrate {
		err = s.autoMigrate()
		if err != nil {
			return nil, fmt.Errorf("failed to autoMigrate: %v", err)
		}
	}

	return s, nil
}

func connect(ctx context.Context, opts *Options) (*gorm.DB, *sql.DB, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("sql store startup deadline exceeded")
		default:
			// db, err := gorm.Open(opts.DatabaseType, opts.URI)
			var dialector gorm.Dialector
			if opts.DatabaseType == DatabaseTypePostgres {
				dialector = postgres.Open(opts.URI)
			} else {
				dialector = sqlite.Open(opts.URI)
			}

			db, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				time.Sleep(1 * time.Second)
				opts.Logger.Warn("sql store connector: can't reach DB, waiting", zap.Error(err))
				continue
			}

			sqlDB, err := db.DB()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get SQL DB: %v", err)
			}

			sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
			sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

			// success
			return db, sqlDB, nil
		}
	}
}

func (s *SQL) autoMigrate() error {
	err := s.db.AutoMigrate(
		&types.Abc{},
	)
	if err != nil {
		s.logger.Error("database migration failed", zap.Error(err))
		return err
	}

	if s.IsDbTypePostgres() {
		// Add FK if needed
	}

	return nil
}

func (s *SQL) IsDbTypePostgres() bool {
	return s.db.Dialector.Name() == "postgres"
}

// Close - closes database connection
func (s *SQL) Close() error {
	return s.sqlDB.Close()
}
