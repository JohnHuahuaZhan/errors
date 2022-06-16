package errors

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func a() {
	fmt.Println("a")
	b()
}
func b() {
	fmt.Println("b")
	c()
}
func c() {
	st := callers(0).StackTrace()
	fmt.Printf("%v\n", st)
	fmt.Printf("%#v", st)
}
func TestCallers(t *testing.T) {
	a()
}
func TestWithStackTrace(t *testing.T) {
	e := NoWrapper("err1")
	e2 := WithMessage(e, "with message")
	e3 := WithStackTrace(e2)
	fmt.Printf("%v", e3)
	fmt.Printf("%+v", e3)
	fmt.Printf("%s", e3)
	fmt.Printf("%+s", e3)
	fmt.Printf("%q", e3)
}
func TestWithStackTrace1(t *testing.T) {
	e := NoWrapper("err1")
	e2 := WithMessage(e, "with message")
	e3 := WithStackTrace(e2)
	e4 := WithStackTrace(e3)
	fmt.Printf("%v\n", e4)
	fmt.Println("----------------------------------------")
	fmt.Printf("%+v\n", e4)
	fmt.Println("----------------------------------------")
	fmt.Printf("%s\n", e4)
	fmt.Println("----------------------------------------")
	fmt.Printf("%+s\n", e4)
	fmt.Println("----------------------------------------")
	fmt.Printf("%q\n", e4)
}
func TestFundamental(t *testing.T) {
	e := New("TestFundamental")
	fmt.Printf("%v\n", e)
	fmt.Println("----------------------------------------")
	fmt.Printf("%s\n", e)
	fmt.Println("----------------------------------------")
	fmt.Printf("%q\n", e)
}
func TestWithCode(t *testing.T) {
	e := New("TestFundamental")
	wc := WithCode(e, 501)
	fmt.Printf("%v\n", wc)
	fmt.Println("----------------------------------------")
	fmt.Printf("%s\n", wc)
	fmt.Println("----------------------------------------")
	fmt.Printf("%q\n", wc)
	fmt.Println("----------------------------------------")
	fmt.Printf("%+v\n", wc)
	fmt.Println("----------------------------------------")
	fmt.Printf("%+s\n", wc)
	fmt.Println("----------------------------------------")
	fmt.Printf("%q\n", wc)
}
func TestWithStackTraceJson(t *testing.T) {
	e := NoWrapper("err1")
	//e2 := WithMessage(e, "with message")
	//e3 := WithStackTrace(e2)
	//e4 := WithStackTrace(e3)
	//e5 := WithMessage(e4, "with message2")
	//e6 := WithCode(e5, 3156)
	//fmt.Printf("%#v\n", e3)
	//fmt.Printf("[%+#v]\n", e3)
	//fmt.Printf("[%#s]\n", e3)
	//fmt.Printf("[%+#s]\n", e3)
	//fmt.Printf("[%+#v]\n", e6)

	e = New("stack")
	for i := 0; i < 200; i++ {
		e = WithMessage(e, strconv.Itoa(i))
	}
	e = WithCode(e, 200)
	fmt.Fprintf(os.Stdout, "%+v", e)
	fmt.Fprintf(os.Stdout, "%#v\n", e)
	fmt.Fprintf(os.Stdout, "%#+v\n", e)
}
