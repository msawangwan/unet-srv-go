package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// deprecate these
const (
	PREFIX_DEBUG = "[DEBUG] "

	PREFIX_ACTIVITY = "[ACTIVITY] "

	PREFIX_INIT       = "[INIT] "
	PREFIX_INIT_ERR   = "[INIT ERROR] "
	PREFIX_INIT_FATAL = "[INIT FATAL] "

	PREFIX_INFO       = "[INFO] "
	PREFIX_INFO_ERR   = "[INFO ERROR] "
	PREFIX_INFO_FATAL = "[INFO FATAL] "

	PREFIX_RESOURCE_ACC     = "[RESOURCE ACCESSED] "
	PREFIX_RESOURCE_ERR     = "[RESOURCE DENIED] "
	PREFIX_RESOURCE_INVALID = "[RESOURCE INVALID] "

	PREFIX_HTTP     = "[HTTP ACTIVITY] "
	PREFIX_HTTP_ERR = "[HTTP ACTIVITY ERROR] "

	PREFIX_DB     = "[DB ACTIVITY] "
	PREFIX_DB_ERR = "[DB ACTIVITY ERROR] "
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
					PREFIX_DEBUG,
					log.Ldate|log.Ltime,
				),
			},
			nil
	}
}

const (
	color_reset = "\033[39m"
	color_red   = "\033[31m"
	color_green = "\033[32m"
)

const (
	colw = 20 // TODO: default column width, from config
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

	l.SetPrefix(fmt.Sprintf(consoleStyle, color_red, ps, color_reset))
}

func (l *Log) PrefixSetWidth(w int, p ...string) {
	consoleStyle = setStyleWidth(w)
	l.Prefix(p...)
}

func (l *Log) PrefixReset() {
	l.SetPrefix(fmt.Sprintf(consoleStyle, color_reset, "DEBUG", color_reset))
}

// DEPRECATE
func (l *Log) SetPrefixDefault() {
	l.SetPrefix(PREFIX_DEBUG)
}

// DEPRECATE
func (l *Log) SetPrefixInit() {
	l.SetPrefix(PREFIX_INIT)
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
