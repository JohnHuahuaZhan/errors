package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Frame uintptr

func (f Frame) pc() uintptr { return uintptr(f) }
func (f Frame) file() string {
	if fu := runtime.FuncForPC(f.pc()); nil != fu {
		file, _ := fu.FileLine(f.pc())
		return file
	}
	return "unknown"
}
func (f Frame) line() int {
	if fu := runtime.FuncForPC(f.pc()); nil != fu {
		_, line := fu.FileLine(f.pc())
		return line
	}
	return -1
}
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %f    源文件名加后缀
//    %+f   完整源文件路径
//    %d    行
//    %n    函数名
//    %+n   带包前缀的函数名
//    %s    同 %f:%d
//    %v    同 %+f:%+n#%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 'f':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
		default:
			io.WriteString(s, onlyFuncName(f.name()))
		}
	case 's':
		fmt.Fprintf(s, "%f", f)
		io.WriteString(s, ":")
		f.Format(s, 'd')
	case 'v':
		fmt.Fprintf(s, "%+f", f)
		io.WriteString(s, ":")
		fmt.Fprintf(s, "%+n", f)
		io.WriteString(s, "#")
		f.Format(s, 'd')
	}
}

type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	用换行组织每一帧，每一帧的内容为Frame的%s
//    %v	用换行组织每一帧，每一帧的内容为Frame的%v
//    %.n   控制显示的stack层数，默认全部
func (st StackTrace) Format(s fmt.State, verb rune) {
	p, _ := s.Precision()

	if p < 1 || p > len(st) {
		p = len(st)
	}
	switch verb {
	case 'v', 's':
		for i := 0; i < p; i++ {
			io.WriteString(s, "\n")
			st[i].Format(s, verb)
		}
	}
}

type stack []uintptr

func (s stack) StackTrace() StackTrace {
	fs := runtime.CallersFrames(s)
	f := make([]Frame, 0, len(s))
	for frame, ok := fs.Next(); ok || (runtime.Frame{} != frame); frame, ok = fs.Next() {
		f = append(f, Frame(frame.PC))
	}
	return f
}

func callers(skip int) stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return st
}

// onlyFuncName removes the path prefix component of a function's name reported by func.Name().
func onlyFuncName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
