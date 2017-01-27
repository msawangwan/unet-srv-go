package adt

import "testing"

var stringtable = []string{
	"alpha",
	"bravo",
	"charlie",
	"delta",
	"fox",
	"echo",
}

var testtable = []struct {
	input    string
	memcheck string
}{
	{"alpha", "alpha"},
	{"bravo", "aslkdjskd"},
}

func TestCanWeAddToTheSet(t *testing.T) {
	t.Log("add members")

	ss := NewStringSet()

	for i, s := range stringtable {
		t.Logf("attempt %d, try to add %s", i, s)
		b := ss.Sadd(s)
		if !b {
			t.Errorf("failed to add [%s]", s)
		} else {
			t.Logf("added %s on attempt %d", s, i)
		}
	}

	t.Logf("%s", ss.ListMem())

	t.Log("test complete")
}

func TestAddTheSameMemberMultipleTimes(t *testing.T) {
	t.Log("a test to see if we can add the same member more than once")

	ss := NewStringSet()

	t.Logf("%s", ss.ListMem())

	for i, s := range testtable {
		t.Logf("on member [%s]", s)
		if !ss.IsMember(s.input) {
			t.Logf("added [%d][%s]", i, s.input)
			ss.Sadd(s.input)
		} else {
			t.Logf("already added [%d][%s]", i, s.input)
		}
		if ss.IsMember(s.input) {
			t.Logf("try to add [%d][%s] a second time", i, s.input)
			try := ss.Sadd(s.input)
			if try {
				t.Errorf("added a duplicate [%d][%s]", i, s.input)
			}
		} else {
			t.Logf("did not add [%d][%s], already a memvber", i, s.input)
		}
	}
	t.Logf("%s", ss.ListMem())

	t.Log("test complete")
}
