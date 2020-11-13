package config

// database configuration (normally we would use postgres)
type Database struct {
	TZ           string `required:"true" envconfig:"TZ"`
	DatabaseType string `default:"sqlite3"` // changes to postgres if details are set
	Sqlite3DB    string `default:"main-app.db"`

	PostgresHost               string `envconfig:"POSTGRES_HOST"`
	PostgresUser               string `envconfig:"POSTGRES_USER"`
	PostgresPassword           string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDB                 string `envconfig:"POSTGRES_DB"`
	PostgresMaxIdleConnections int    `default:"100" envconfig:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	MaxOpenConnections         int    `default:"40" envconfig:"MAX_OPEN_DB_CONNECTIONS"`
}

type Logger struct {
	DisableStackTrace    bool `default:"false" envconfig:"LOGGER_DISABLE_STACK_TRACE"`
	ProductionFormatting bool `default:"true" envconfig:"LOGGER_PRODUCTION_FORMATTING"`
}

type MainAppServer struct {
	EnableCors bool   `default:"false" envconfig:"MAIN_APP_ENABLE_CORS"`
	HttpHost   string `default:"localhost" envconfig:"MAIN_APP_HTTP_HOST"`
	HttpPort   int    `default:"9701" envconfig:"MAIN_APP_SERVER_HTTP_PORT"`
}

type (
	// Main app configurations
	MainAppConfig struct {
		Database
		Logger
		MainAppServer
	}

	// Main app configurations
	MainAppCliConfig struct {
		Logger
	}
)
