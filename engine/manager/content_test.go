package manager

import "testing"

func TestDoesTheContentManagerHandler(t *testing.T) {
	t.Log("does the content handler work?")

	ch, e := NewContentHandler(nil)
	if e != nil {
		t.Errorf("%s", e)
	}

	r, e := ch.Random()
	if e != nil {
		t.Errorf("%s", e)
	}

	n := ch.GenerateHyphenatedName()

	t.Logf("some lore: %s and a name: %s", *r, n)
}
