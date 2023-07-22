// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2021-12-19
// Last changed: 2023-07-22 03:01:47
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU Affero General Public License
// as published by the Free Software Foundation, either version 3 of
// the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program. If not, see
// <http://www.gnu.org/licenses/>.

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
