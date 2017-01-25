package data

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"encoding/json"
)

type Name string
type Label string
type Type string
type Category string

type Properties []string

type PlayerOwned bool

func (po PlayerOwned) String() string {
	var phrase string = "player owned: "
	if po {
		return fmt.Sprintf("%s%s", phrase, "true")
	}
	return fmt.Sprintf("%s%s", phrase, "false")
}

type Schema struct {
	DataLabel       Label       `json:"label"`
	DataType        Type        `json:"type"`
	DataCategory    Category    `json:"category"`
	DataProperties  Properties  `json:"properties"`
	DataPlayerOwned PlayerOwned `json:"playerOwned"`
}

type Schemas map[Name]Schema

type SchemaCache struct {
	Schemas `json:"schemas"`
}

func (sc *SchemaCache) PrettyPrint() string {
	var s string = ""
	put := func(str interface{}) { s = s + fmt.Sprintf("%v\n", str) }
	for k, v := range sc.Schemas {
		put("*** " + k + " ***")
		put("\t" + "label: " + v.DataLabel)
		put("\t" + "type: " + v.DataType)
		put("\t" + "cat: " + v.DataCategory)
		put("\t" + v.DataPlayerOwned.String())
		put("\t" + "properties: ")
		for _, e := range v.DataProperties {
			put("\t\t" + e)
		}
	}
	return s
}

func MarshallSchema(fn string) (*SchemaCache, error) {
	if !strings.HasSuffix(fn, ".json") {
		return nil, errors.New("data file not a json file")
	}

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s SchemaCache

	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func debugPrintDataFileByLine(fn string) error {
	if !strings.HasSuffix(fn, ".json") {
		return errors.New("data file not a json file")
	}

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	if err = s.Err(); err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}
