package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type level int

const (
	NONE level = iota
	VERBOSE
	DEBUG
	WARN
	ERROR
	FATAL
)

type Log struct {
	*log.Logger
}

func NewLogger(filename string) (*Log, error) {
	logfile, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, err
	} else {
		return &Log{
				Logger: log.New(
					io.MultiWriter(logfile, os.Stdout),
					"CONSOLE",
					log.Ldate|log.Ltime,
				),
			},
			nil
	}
}

const (
	color_black   = "\033[30m"
	color_red     = "\033[31m"
	color_green   = "\033[32m"
	color_yellow  = "\033[33m"
	color_blue    = "\033[34m"
	color_magenta = "\033[35m"
	color_cyan    = "\033[36m"
	color_gray    = "\033[37m"
	color_white   = "\033[97m"
	color_reset   = "\033[39m"
)

const (
	colw = 40 // TODO: default column width, from config
)

var (
	consoleStyle = setStyleWidth(colw)
)

// TODO: DEPRECATED
func (l *Log) Prefix(p ...string) {
	var (
		ps string
	)

	for _, pf := range p {
		ps = ps + "|" + strings.ToUpper(pf) + "|"
	}

	l.SetPrefix(fmt.Sprintf(consoleStyle, color_blue, ps, color_reset))
}

// TODO: DEPRECATED
func (l *Log) PrefixError(p ...string) {
	var (
		ps string
	)

	for _, pf := range p {
		ps = ps + "|" + strings.ToUpper(pf) + "|"
	}

	l.SetPrefix(fmt.Sprintf(consoleStyle, color_red, ps, color_reset))
}

// TODO: DEPRECATE
func (l *Log) PrefixReset() {
	l.SetPrefix(fmt.Sprintf(consoleStyle, color_reset, "DEBUG", color_reset))
}

func (l *Log) Label(lvl level, p ...string) {
	var ps, color string

	for _, pf := range p {
		ps = ps + "|" + strings.ToUpper(pf) + "|"
	}

	switch lvl {
	case NONE:
		color = color_gray
	case VERBOSE:
		color = color_cyan
	case DEBUG:
		color = color_green
	case WARN:
		color = color_yellow
	case ERROR:
		color = color_magenta
	case FATAL:
		color = color_red
	default:
		color = color_white
	}

	l.SetPrefix(fmt.Sprintf(consoleStyle, color, ps, color_reset))
}

func (l *Log) ClearLabel() {
	l.SetPrefix(fmt.Sprintf(consoleStyle, color_reset, "DEBUG", color_reset))
}

func (l *Log) Level(lvl level) {
	switch lvl {
	case DEBUG:
		l.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	case VERBOSE:
		l.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	default:
		l.SetFlags(log.Ldate | log.Ltime)
	}
}

func (l *Log) PrefixSetWidth(w int, lvl level, p ...string) {
	consoleStyle = setStyleWidth(w)
	l.Label(lvl, p...)
}

func setStyleWidth(w int) string {
	return "%s%-" + strconv.Itoa(w) + "s%s "
}
