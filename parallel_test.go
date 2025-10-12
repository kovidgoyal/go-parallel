package parallel

import (
	"fmt"
	"strings"
	"testing"
)

var _ = fmt.Print

func panicer(r any) *PanicError {
	return Format_stacktrace_on_panic(r, 0)
}

func panicking_callback(start, limit int) {
	panic("panicking_callback")
}

func TestStacktrace(t *testing.T) {
	a := panicer("XXX")
	if !strings.HasSuffix(a.frames[0].Function, ".panicer") {
		t.Fatalf("Most recent call is %s instead of panicer()", a.frames[0].Function)
	}
	b := panicer(a)
	text := b.Error()
	if strings.Count(text, "Stack trace") != 2 {
		t.Fatalf("There are not two stack traces in:\n%s", text)
	}

	err := Run_in_parallel_over_range(1, panicking_callback, 0, 100)
	a = err.(*PanicError)
	if !strings.HasSuffix(a.frames[0].Function, ".panicking_callback") {
		t.Fatalf("Most recent call is %s instead of panicer()", a.frames[0].Function)
	}
}
