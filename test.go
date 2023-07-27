package nostr

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(inRed("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(inRed("\tcatch: %s\n\n"), err.Error())
		fmt.Printf(out)
		tb.Fail()
	}
}

// Fails the test if "want" is not equal to "got".
func equals(tb testing.TB, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(inRed("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(inRed("\twant: %#v\n\n"), want)
		out += fmt.Sprintf(inRed("\tgot: %#v\n\n"), got)
		fmt.Printf(out)
		tb.Fail()
	}
}

// Fails the test if the given file doesn't exist.
func exists(tb testing.TB, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, file, line, _ := runtime.Caller(1)
		out := fmt.Sprintf(inRed("%s:%d\n\n"), filepath.Base(file), line)
		out += fmt.Sprintf(inRed("\tno file: %s\n\n"), filename)
		fmt.Printf(out)
		tb.Fail()
	}
}
