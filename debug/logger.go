package debug

import (
	"io"
	"log"
	"os"
)

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

func (l *Log) SetPrefixDefault() {
	l.SetPrefix(PREFIX_DEBUG)
}

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
