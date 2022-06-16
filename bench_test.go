package errors

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

func Errors(at, depth int) error {
	if at >= depth {
		return New("no error")
	}
	return Errors(at+1, depth)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE interface{}

func BenchmarkErrors(b *testing.B) {
	type run struct {
		stack int
		pkg   string
	}
	runs := []run{
		{10, "std"},
		{100, "std"},
		{1000, "std"},
	}
	for _, r := range runs {
		var part string
		var f func(at, depth int) error
		switch r.pkg {
		case "std":
			part = "errors"
			f = Errors
		default:
		}

		name := fmt.Sprintf("%s-stack-%d", part, r.stack)
		b.Run(name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = f(0, r.stack)
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}
func BenchmarkFormat(b *testing.B) {
	err := New("great")
	b.Run("NORMALS", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%s", err)
		}
		b.StopTimer()
	})
	b.Run("NORMALV", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%v", err)
		}
		b.StopTimer()
	})
	err1 := WithMessage(err, "ok")
	err2 := WithCode(err1, 404)

	b.Run("JsonS", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%#s", err2)
		}
		b.StopTimer()
	})
	b.Run("JsonV", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%#v", err2)
		}
		b.StopTimer()
	})
	b.Run("LookBackS", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%+s", err2)
		}
		b.StopTimer()
	})
	b.Run("LookBackV", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%+v", err2)
		}
		b.StopTimer()
	})
	b.Run("LookBackSJson", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%+#s", err2)
		}
		b.StopTimer()
	})
	b.Run("LookBackVJson", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(io.Discard, "%+#v", err2)
		}
		b.StopTimer()
	})

	b.Run("NewAndLookBackVJson", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			e := New("great")
			fmt.Fprintf(io.Discard, "%+#v", e)
		}
		b.StopTimer()
	})

}
func BenchmarkDeepLookBackVJson(b *testing.B) {
	e := New("stack")
	for i := 0; i < 200; i++ {
		e = WithMessage(e, strconv.Itoa(i))
	}
	e = WithCode(e, 200)
	fmt.Fprintf(os.Stdout, "[%+#v]", e)
	b.Run("NewAndLookBackVJson", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			e := New("great")
			fmt.Fprintf(io.Discard, "%+#v", e)
		}
		b.StopTimer()
	})
}
