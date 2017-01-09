package handler

import (
	"fmt"
	"io"
	"strconv"

	"encoding/json"

	"github.com/msawangwan/unet-srv-go/debug"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

func raiseServerError(err error) exception.Handler {
	return raise(err, err.Error(), 500)
}

func setPrefix(prefix string, specific string, l *debug.Log) func() {
	l.SetPrefix(fmt.Sprintf("[HTTP][HANDLER][%s][%s] ", prefix, specific))
	return func() {
		defer l.SetPrefixDefault()
	}
}

func parseJSON(r io.Reader) (interface{}, error) {
	// setPrefix("JSON", "PARSE", )
	var (
		j interface{}
	)

	err := json.NewDecoder(r).Decode(&j)

	if err != nil {
		return nil, err
	}

	return j.(map[string]interface{}), nil
}

func parseJSONInt(r io.Reader) (*int, error) {
	var (
		j interface{}
	)

	d := json.NewDecoder(r)
	d.UseNumber()
	if err := d.Decode(&j); err != nil {
		return nil, err
	}

	n := j.(map[string]interface{})["value"].(json.Number)
	val, _ := strconv.Atoi(string(n))
	return &val, nil
}
