package debug

import (
	"fmt"
	"testing"
)

func TestPutColors(t *testing.T) {
	t.Log("print colors ...")
	var s = []string{
		fmt.Sprintf("%sred%s", color_red, color_reset),
		fmt.Sprintf("%sgreen%s", color_green, color_reset),
		fmt.Sprintf("%syellow%s", color_yellow, color_reset),
		fmt.Sprintf("%sblue%s", color_blue, color_reset),
		fmt.Sprintf("%swhite%s", color_white, color_reset),
	}

	for _, v := range s {
		t.Log(v)
	}
}
