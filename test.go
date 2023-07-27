package nostr

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var (
	reset = "\033[0m"
	red   = "\033[31m"
)

func colorize(color, s string) string {
	return color + s + reset
}

func redf(s string) string {
	return colorize(red, s)
}

// Fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(redf("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(redf("\tcatch: %s\n\n"), err.Error())
		fmt.Printf(out)
		tb.Fail()
	}
}

// Fails the test if "want" is not equal to "got".
func equals(tb testing.TB, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(redf("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(redf("\twant: %#v\n\n"), want)
		out += fmt.Sprintf(redf("\tgot: %#v\n\n"), got)
		fmt.Printf(out)
		tb.Fail()
	}
}

// Fails the test if the given file doesn't exist.
func exists(tb testing.TB, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(redf("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(redf("\tno file: %s\n\n"), filename)
		fmt.Printf(out)
		tb.Fail()
	}
}
