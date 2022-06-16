// v的内容要丰富许多 使用了Frame的flag+,s则未使用
// 除了Frame，其他类型的flag+只用于设置是否回溯
// flag# 可将各个错误用json格式输出，如果配合+,那么这些json将会被,连接起来。请在打印前后加上[]组成一个json数组

package errors

import (
	"fmt"
	"io"
	"strconv"
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
	if cause == nil {
		return nil
	}
	e := &withMessage{
		cause: cause,
		msg:   msg,
	}
	return e
}
func (nw *noWrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		if s.Flag('#') {
			fmt.Fprintf(s, "{\"msg\":%s}", strconv.Quote(nw.Error()))
		} else {
			io.WriteString(s, nw.Error())
			io.WriteString(s, "\n")
		}

	case 'q':
		io.WriteString(s, nw.Error())
	}
}

func WithMessageF(cause error, format string, args ...interface{}) error {
	return WithMessage(cause, fmt.Sprintf(format, args...))
}
func (wm *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		if s.Flag('+') {
			format := fmt.Sprintf("%%%s", string(verb))
			lookBack(s, verb, format, wm.Unwrap())
		}
		if s.Flag('#') {
			if s.Flag('+') {
				io.WriteString(s, ",")
			}
			fmt.Fprintf(s, "{\"msg\":%s}", strconv.Quote(wm.Error()))
		} else {
			fmt.Fprintf(s, "%s\n", wm.Error())
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
	case 'v', 's':
		if s.Flag('+') {
			format := fmt.Sprintf("%%%s", string(verb))
			lookBack(s, verb, format, w.Unwrap())
		}
		if s.Flag('#') {
			if s.Flag('+') {
				io.WriteString(s, ",")
			}
			io.WriteString(s, "{\"trace\":")
			w.trace.Format(s, verb)
			io.WriteString(s, "}")
		} else {
			w.trace.Format(s, verb)
		}
	case 'q': //兼容q%
		fmt.Fprintf(s, "%q", w.Error())
	}

}

func lookBack(s fmt.State, verb rune, format string, err error) {
	f, ok := err.(fmt.Formatter)
	if ok {
		f.Format(s, verb)
	} else {
		fmt.Fprintf(s, format, err)
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
		if s.Flag('#') {
			fmt.Fprintf(s, "{\"msg\":%s,\"trace\":", strconv.Quote(f.Error()))
			f.trace.Format(s, verb)
			io.WriteString(s, "}")
		} else {
			io.WriteString(s, f.Error())
			io.WriteString(s, "\n")
			f.trace.Format(s, verb)
		}
	case 'q':
		io.WriteString(s, f.Error())
	}

}

func WithCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &withCode{
		error: err,
		code:  code,
	}
}

type withCode struct {
	error
	code int
}

func (wc *withCode) Unwrap() error {
	return wc.error
}

func (wc *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		if s.Flag('+') {
			format := fmt.Sprintf("%%%s", string(verb))
			lookBack(s, verb, format, wc.Unwrap())
		}
		if s.Flag('#') {
			if s.Flag('+') {
				io.WriteString(s, ",")
			}
			fmt.Fprintf(s, "{\"code\":%d}", wc.code)
		} else {
			fmt.Fprintf(s, "code:%d\n", wc.code)
		}

	case 'q':
		io.WriteString(s, wc.Error())
	}
}
