package debug

/*
 *
 * see: http://nickgravgaard.com/elastic-tabstops/index.html
 * and also see: https://golang.org/pkg/text/tabwriter/
 *
 * this is nearly the same thing i came up with:
 * https://github.com/siongui/userpages/blob/master/content/code/go-fixed-width-string/fixedwidth.go
 *
 * colors: https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
 * ncurses: https://en.wikipedia.org/wiki/Ncurses
 *
 * finally some good stuff here:
 * http://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/
 *
 * this post shows how color codes are used in go:
 * http://stackoverflow.com/questions/27242652/colorizing-golang-test-run-output
 *
 */

import (
	"fmt"
)

//const (
//	color_reset = "\033[39m"
//	color_red   = "\033[31m"
//	color_green = "\033[32m"
//)

const (
	style1 = "%40s |" + color_red + "%40s\033[39m"
	style2 = "%5s |" + color_green + "%40s\033[39m"
	style3 = "%5s" + color_green + "%40s\033[39m"
)

type Printer interface {
	Put() string
}

type ConsoleFormatter struct{}

func (cf *ConsoleFormatter) Put(text string) string {
	return fmt.Sprintf(style3, "", text)
}
