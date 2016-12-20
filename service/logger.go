package service

import (
	//	"fmt"
	//	"io"
	"log"
	//	"net/http"
	//	"os"
)

const (
	PREFIX_DEFAULT = "[ACTIVITY] "

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

type ServiceLogger interface {
	Log(s string)
	Init(s string)
	Info(s string)
	Resource(s string)
	Http(s string)
	Db(s string)
}

type Console struct {
	*log.Logger
}

func (c *Console) Log(s string)      {}
func (c *Console) Init(s string)     {}
func (c *Console) Info(s string)     {}
func (c *Console) Resource(s string) {}
func (c *Console) Http(s string)     {}
func (c *Console) Db(s string)       {}

type Empty struct{}

func (nl *Empty) Log(s string)      {}
func (nl *Empty) Init(s string)     {}
func (nl *Empty) Info(s string)     {}
func (nl *Empty) Resource(s string) {}
func (nl *Empty) Http(s string)     {}
func (nl *Empty) Db(s string)       {}
