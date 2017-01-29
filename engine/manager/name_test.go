package manager

import "testing"

func TestGenerateOneName(t *testing.T) {
	t.Log("lets print a single name")

	ng, e := NewNameGenerator()
	if e != nil {
		t.Errorf("error creating the name generator")
	}

	name := ng.GenerateHyphenatedName()
	t.Logf("generated a name ... [%s]", name)

	t.Log("test complete")
}

func TestGenerateMultipleName(t *testing.T) {
	t.Log("lets print some names")

	ng, e := NewNameGenerator()
	if e != nil {
		t.Errorf("error creating the name generator")
	}

	var i, n = 0, 50

	for i < n {
		name := ng.GenerateHyphenatedName()
		t.Logf("generated a name ... [%s]", name)
		i++
	}

	t.Log("test complete")
}
