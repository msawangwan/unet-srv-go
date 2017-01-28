package serialization

import "testing"

type JSONTestObjectOne struct {
	ObjOne_FieldOne string `json:"objOne_fieldOne"`
	ObjOne_FieldTwo string `json:"objOne_fieldTwo"`
}

type JSONTestObjectTwo struct {
	ObjTwo_FieldOne string
	ObjTwo_FieldTwo string
}

func TestCanWeLoadGenericTypesFromJSON(t *testing.T) {
	t.Log("can we load generic types from json")

	var testobject JSONTestObjectOne

	testobjinterface, e := LoadJSONFile("test.json", testobject)
	if e != nil {
		t.Errorf("%s", e)
	}

	o2 := testobjinterface.(JSONTestObjectOne)

	t.Logf("%v", o2)
}
