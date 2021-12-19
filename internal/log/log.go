package log

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	logWr := os.Stderr
	isTerm := isatty.IsTerminal(logWr.Fd())

	consoleWriter := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339,
		Out:        os.Stdout,
		NoColor:    !isTerm}
	log = zerolog.New(consoleWriter).With().Timestamp().Logger()
	log = log.Level(zerolog.WarnLevel)
}

func SetLevel(lvl int) {
	switch {
	case lvl < 1:
		return
	case lvl == 1:
		log = log.Level(zerolog.InfoLevel)
	case lvl == 2:
		log = log.Level(zerolog.DebugLevel)
	case lvl > 2:
		log = log.Level(zerolog.TraceLevel)
	}
}

func Init(dir string) {
	logWr := os.Stderr
	isTerm := isatty.IsTerminal(logWr.Fd())

	consoleWriter := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339,
		Out:        os.Stdout,
		NoColor:    !isTerm}

	f, _ := os.Create(fmt.Sprintf("%s/app.log", dir))

	fileWriter := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339,
		Out:        f,
		NoColor:    true}

	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	log = zerolog.New(multi).With().Timestamp().Logger()
}

func Info(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
func Debug(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}
func Warn(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}
func Error(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}
