package logger

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/suhriar/blog-mono-api/config"
)

// InitializeLogger mengatur konfigurasi logger global
func InitializeLogger(cfg *config.Config) func() {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(level)

	var stdOut io.Writer = os.Stdout
	if cfg.Log.Type == "text" {
		stdOut = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}
	writers := []io.Writer{stdOut}
	var runLogFile *os.File
	if cfg.Log.LogFileEnabled {
		runLogFile, err = os.OpenFile(
			cfg.Log.LogFilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0666,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to open log file")
		}

		writers = append(writers, runLogFile)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return func() {
		if runLogFile != nil {
			runLogFile.Close()
		}
	}
}

func RequestLogger(r *http.Request) zerolog.Logger {
	logger := zerolog.Ctx(r.Context())
	if logger.GetLevel() == zerolog.Disabled {
		return log.Logger
	}
	return *logger
}
