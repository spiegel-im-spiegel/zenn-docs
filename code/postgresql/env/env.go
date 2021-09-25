package env

import (
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/gocli/config"
)

const (
	ServiceName = "elephantsql"
)

func init() {
	//load ${XDG_CONFIG_HOME}/${ServiceName}/env file
	if err := godotenv.Load(config.Path(ServiceName, "env")); err != nil {
		//load .env file
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}
}

func PostgresDSN() string {
	return os.Getenv("ELEPHANTSQL_URL")
}

func LogLevel() LoggerLevel {
	return getLogLevel(os.Getenv("LOGLEVEL"))
}

func ZerologLevel() zerolog.Level {
	return LogLevel().ZerlogLevel()
}

func PgxlogLevel() pgx.LogLevel {
	return LogLevel().PgxLogLevel()
}

func EnableLogFile() bool {
	return strings.EqualFold(os.Getenv("ENABLE_LOGFILE"), "true")
}
