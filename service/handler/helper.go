package handler

import (
	"errors"
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
		l.SetPrefixDefault()
	}
}

func parseJSON(r io.Reader) (interface{}, error) {
	var (
		j interface{}
	)

	err := json.NewDecoder(r).Decode(&j)

	if err != nil {
		return nil, err
	}

	return j, nil
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

func marshallJSONString(j interface{}) (*string, error) {
	var (
		s string
	)

	m, ok := j.(map[string]interface{})
	if ok {
		s, ok = m["value"].(string)
		if !ok {
			return nil, errors.New("failed to read json map (line 58)")
		}
	} else {
		return nil, errors.New("failed to parse json struct (line 67)")
	}

	return *s, nil
}
