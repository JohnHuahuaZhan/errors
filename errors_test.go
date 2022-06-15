package errors

import (
	"fmt"
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
	fmt.Printf("%+a", st)
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
