package manager

import "testing"

const (
	path = "../../lore.json"
)

func TestCanWeParseTheLoreTextFile(t *testing.T) {
	lg, e := NewLoreGenerator(path, nil)
	if e != nil {
		t.Errorf("%s", e)
	}

	for i, v := range lg.cache {
		t.Logf("[key: %s] [value: %s]", i, v)
	}

	t.Log("test complete")
}

func TestCanWeGetARandomDescription(t *testing.T) {
	lg, e := NewLoreGenerator(path, nil)
	if e != nil {
		t.Errorf("%s", e)
	}

	var i = 0
	for i < 50 {
		r, e := lg.Random()
		if e != nil {
			t.Errorf("%s", e)
		}

		t.Logf("picking a random description ...[%s]", *r)
		i++
	}

	t.Log("test complete")
}
