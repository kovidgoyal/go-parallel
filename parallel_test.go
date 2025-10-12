package parallel

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var _ = fmt.Print

func panicer(r any) *PanicError {
	return Format_stacktrace_on_panic(r, 0)
}

func panicking_callback(start, limit int) {
	panic("panicking_callback")
}

func filler(items []int, start, limit int) {
	for i := start; i < limit; i++ {
		items[i] = i
	}
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

	expected := make([]int, 0, 64)
	filler(expected, 0, len(expected))
	for num := range 4 {
		err := Run_in_parallel_over_range(num, panicking_callback, 0, 100)
		a = err.(*PanicError)
		if !strings.HasSuffix(a.frames[0].Function, ".panicking_callback") {
			t.Fatalf("Most recent call is %s instead of panicer()", a.frames[0].Function)
		}

		items := make([]int, 0, len(expected))
		if err = Run_in_parallel_over_range(num, func(start, limit int) { filler(items, start, limit) }, 0, len(items)); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, items); diff != "" {
			t.Fatalf("Not filled for num=%d\n%s", num, diff)
		}
	}
}
