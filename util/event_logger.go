package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	LOGFILE_NAME           = "access.log"
	PREFIX_FATAL_ERR       = "[FATAL] "
	PREFIX_DEFAULT         = "[EVENT] "
	PREFIX_INIT            = "[INIT] "
	PREFIX_INIT_FATAL      = "[INIT FATAL] "
	PREFIX_RESOURCE_ACCESS = "[RESOURCE ACCESSED] "
	PREFIX_INVALID_REQUEST = "[INVALID RESOURCE REQUEST] "
	PREFIX_HTTP_ERR        = "[HTTP ERROR] "
	PREFIX_DB_ACTIVITY     = "[DB ACTIVITY] "
	PREFIX_DB_ERR          = "[DB ERROR] "
)

type eventLogger struct {
	event *log.Logger
}

var Log *eventLogger

func init() {
	logfile, err := os.OpenFile(
		LOGFILE_NAME,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		log.Printf("error opening %s: %v\n", LOGFILE_NAME, err)
	}

	Log = &eventLogger{
		event: log.New(
			io.MultiWriter(logfile, os.Stdout),
			PREFIX_DEFAULT,
			log.Ldate|log.Ltime,
		),
	}
}

func (el *eventLogger) Fatal(err error) {
	reason := fmt.Sprintf("a fatal error has occured: %s", err.Error())
	el.event.SetPrefix(PREFIX_FATAL_ERR)
	el.event.Fatalf("%s\n", reason)
}

func (el *eventLogger) InitMessage(s string) {
	reason := fmt.Sprintf("setup complete: %s", s)
	el.event.SetPrefix(PREFIX_INIT)
	el.event.Printf("%s\n", reason)
}

func (el *eventLogger) InitPanic(err error) {
	reason := fmt.Sprintf("fatal error on setup: %s", err.Error())
	el.event.SetPrefix(PREFIX_INIT_FATAL)
	el.event.Panic(reason)
}

func (el *eventLogger) ResourceRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reason := fmt.Sprintf("resource request (%s): %s", http.StatusText(200), r.URL.Path)
		el.event.SetPrefix(PREFIX_RESOURCE_ACCESS)
		el.event.Printf("%s\n", reason)
		h(w, r)
	}
}

func (el *eventLogger) DbActivity(s string) {
	el.event.SetPrefix(PREFIX_DB_ACTIVITY)
	el.event.Printf("%s\n", s)
}

func (el *eventLogger) InvalidRequest(w http.ResponseWriter, r *http.Request) {
	reason := fmt.Sprintf("invalid resource request (%s): %s", http.StatusText(500), r.URL.Path)
	http.Error(w, reason, 500)
	el.event.SetPrefix(PREFIX_INVALID_REQUEST)
	el.event.Printf("%s\n", reason)
}

func (el *eventLogger) DbErr(w http.ResponseWriter, r *http.Request, e error) {
	reason := fmt.Sprintf("an error occured while accessing the db %s (%s): %s", r.URL.Path, http.StatusText(500), e.Error())
	el.event.SetPrefix(PREFIX_DB_ERR)
	el.event.Printf("%s\n", reason)
}

func (el *eventLogger) HttpErr(w http.ResponseWriter, r *http.Request, e error) {
	reason := fmt.Sprintf("an error occured while processing the request for %s (%s): %s", r.URL.Path, http.StatusText(500), e.Error())
	el.event.SetPrefix(PREFIX_HTTP_ERR)
	el.event.Printf("%s\n", reason)
}
