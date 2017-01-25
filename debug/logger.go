package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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
	color_reset  = "\033[39m"
	color_red    = "\033[31m"
	color_green  = "\033[32m"
	color_yellow = "\033[33m"
	color_blue   = "\033[34m"
	color_white  = "\033[37m"
)

const (
	colw = 40 // TODO: default column width, from config
)

var (
	consoleStyle = setStyleWidth(colw)
)

func (l *Log) Prefix(p ...string) {
	var (
		ps string
	)

	for _, pf := range p {
		ps = ps + "|" + strings.ToUpper(pf) + "|"
	}

	l.SetPrefix(fmt.Sprintf(consoleStyle, color_blue, ps, color_reset))
}

func (l *Log) PrefixError(p ...string) {
	var (
		ps string
	)

	for _, pf := range p {
		ps = ps + "|" + strings.ToUpper(pf) + "|"
	}

	l.SetPrefix(fmt.Sprintf(consoleStyle, color_red, ps, color_reset))
}

func (l *Log) PrefixSetWidth(w int, p ...string) {
	consoleStyle = setStyleWidth(w)
	l.Prefix(p...)
}

func (l *Log) PrefixReset() {
	l.SetPrefix(fmt.Sprintf(consoleStyle, color_reset, "DEBUG", color_reset))
}

func (l *Log) SetLevelDefault() {
	l.SetFlags(log.Ldate | log.Ltime)
}

func (l *Log) SetLevelDebug() {
	l.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func (l *Log) SetLevelVerbose() {
	l.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func setStyleWidth(w int) string {
	return "%s%-" + strconv.Itoa(w) + "s%s "
}
