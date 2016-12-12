package service

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	LOGFILE_NAME           = "access.log"
	PREFIX_RESOURCE_ACCESS = "[Resource Access Request] "
)

type gatewayLogger struct {
	resourceAccessLog *log.Logger
}

var Log *gatewayLogger

func init() {
	logfile, err := os.OpenFile(
		LOGFILE_NAME,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		log.Printf("error opening %s: %v\n", LOGFILE_NAME, err)
	}

	Log = &gatewayLogger{
		resourceAccessLog: log.New(
			io.MultiWriter(logfile, os.Stdout),
			PREFIX_RESOURCE_ACCESS,
			log.Ldate|log.Ltime,
		),
	}
}

func (gl *gatewayLogger) resourceRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gl.resourceAccessLog.Printf("served %v\n", r.URL.Path)
		h(w, r)
	}
}
