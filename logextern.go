package log

import (
	"fmt"
	"sync/atomic"
)

const (
	FATAL = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
	TRACE = 6
	MAX   = 7
)

type LoggerEx struct {
	Level int
}

func (l *Logger) LogLevel(lv int) int {
	if lv >= MAX || lv < FATAL {
		return l.Level
	}
	l.Level = lv
	return l.Level
}

func LogLevel(lv int) int {
	return std.LogLevel(lv)
}

func Infof(format string, v ...any) {
	if atomic.LoadInt32(&std.isDiscard) != 0 {
		return
	}
	if std.Level < INFO {
		return
	}
	std.Output(2, fmt.Sprintf(format, v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < INFO {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...any) {
	if atomic.LoadInt32(&std.isDiscard) != 0 {
		return
	}
	if std.Level < DEBUG {
		return
	}
	std.Output(2, fmt.Sprintf(format, v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < DEBUG {
		return
	}
	fmt.Printf("lev:%v\n", l.Level)
	l.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...any) {
	if atomic.LoadInt32(&std.isDiscard) != 0 {
		return
	}
	if std.Level < ERROR {
		return
	}
	if g_advLog != nil {
		g_advLog.Output(2, fmt.Sprintf(format, v...))
	}
	std.Output(2, fmt.Sprintf(format, v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < ERROR {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...any) {
	if atomic.LoadInt32(&std.isDiscard) != 0 {
		return
	}
	if std.Level < WARN {
		return
	}
	if g_advLog != nil {
		g_advLog.Output(2, fmt.Sprintf(format, v...))
	}
	std.Output(2, fmt.Sprintf(format, v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < WARN {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Tracef(format string, v ...any) {
	if atomic.LoadInt32(&std.isDiscard) != 0 {
		return
	}
	if std.Level < TRACE {
		return
	}
	std.Output(2, fmt.Sprintf(format, v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Tracef(format string, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < TRACE {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}
