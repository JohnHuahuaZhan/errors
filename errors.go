// v的内容要丰富许多 使用了Frame的flag+,s则未使用
// 除了Frame，其他类型的flag+只用于设置是否回溯

package errors

import (
	"fmt"
	"io"
)

func NoWrapper(msg string) error {
	return &noWrapper{msg: msg}
}
func NoWrapperF(format string, args ...interface{}) error {
	return &noWrapper{msg: fmt.Sprintf(format, args...)}
}

type noWrapper struct {
	msg string
}

func (nw *noWrapper) Error() string {
	return nw.msg
}
func WithMessage(cause error, msg string) error {
	e := &withMessage{
		cause: cause,
		msg:   msg,
	}
	return e
}
func (nw *noWrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's', 'q':
		io.WriteString(s, nw.Error())
	}
}

func WithMessageF(cause error, format string, args ...interface{}) error {
	e := &withMessage{
		cause: cause,
		msg:   fmt.Sprintf(format, args...),
	}
	return e
}
func (wm *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", wm.Unwrap()) //需要+才能一直回溯下去
			io.WriteString(s, wm.msg)
		} else {
			fmt.Fprintf(s, "%v\n", wm.msg)
		}
	case 's':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+s\n", wm.Unwrap())
			io.WriteString(s, wm.msg)
		} else {
			fmt.Fprintf(s, "%s\n", wm.msg)
		}
	case 'q':
		io.WriteString(s, wm.Error())
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (wm *withMessage) Error() string {
	return wm.msg
}
func (wm *withMessage) Unwrap() error {
	return wm.cause
}

//WithStackTrace 堆栈从WithStack调用处开始
func WithStackTrace(err error) error {
	if err == nil {
		return nil
	}
	return &withStackTrace{
		err,
		callers(3).StackTrace(),
	}
}

type withStackTrace struct {
	error
	trace StackTrace
}

func (w *withStackTrace) Unwrap() error { return w.error }

//Format %v 同StackTrace的%v
// %+v 会在%v的基础上向上追溯打印wrapper的内容
// %s同理
func (w *withStackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.trace.Format(s, verb)
		} else {
			w.trace.Format(s, verb)
		}
	case 's':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+s", w.Unwrap())
			w.trace.Format(s, verb)
			return
		} else {
			w.trace.Format(s, verb)
		}
	case 'q': //兼容q%
		fmt.Fprintf(s, "%q", w.Error())
	}

}

func New(message string) error {
	return &fundamental{
		msg:   message,
		trace: callers(3).StackTrace(),
	}
}
func NewF(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		trace: callers(3).StackTrace(),
	}
}

// fundamental no wrapper
type fundamental struct {
	msg   string
	trace StackTrace
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		io.WriteString(s, f.msg)
		f.trace.Format(s, verb)
	case 'q': //兼容q%
		fmt.Fprintf(s, "%q", f.Error())
	}

}
