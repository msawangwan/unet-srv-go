package debug

import "testing"

func TestWhatDoesFormatterDo(t *testing.T) {
	t.Log("test the console formatter")
	cf := &ConsoleFormatter{}

	t.Log(cf.Put("this is some other really long test"))
	t.Log(cf.Put("this is only a test hm"))
	t.Log(cf.Put("this is some other really long test"))
	t.Log(cf.Put("test"))
	t.Log(cf.Put("is only a test"))

	t.Log("test complete")
}
