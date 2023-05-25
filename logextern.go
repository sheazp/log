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

func Infof(format string, v ...interface{}) {
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
func (l *Logger) Infof(format string, v ...interface{}) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < INFO {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
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
func (l *Logger) Debugf(format string, v ...interface{}) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < DEBUG {
		return
	}
	fmt.Printf("lev:%v\n", l.Level)
	l.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
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
func (l *Logger) Errorf(format string, v ...interface{}) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < ERROR {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
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
func (l *Logger) Warnf(format string, v ...interface{}) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < WARN {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

func Tracef(format string, v ...interface{}) {
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
func (l *Logger) Tracef(format string, v ...interface{}) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if l.Level < TRACE {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func Black(msg string) string {
	return SetColor(msg, 0, 0, TextBlack)
}

func Red(msg string) string {
	return SetColor(msg, 0, 0, TextRed)
}

func Green(msg string) string {
	return SetColor(msg, 0, 0, TextGreen)
}

func Yellow(msg string) string {
	return SetColor(msg, 0, 0, TextYellow)
}

func Blue(msg string) string {
	return SetColor(msg, 0, 0, TextBlue)
}

func Magenta(msg string) string {
	return SetColor(msg, 0, 0, TextMagenta)
}

func Cyan(msg string) string {
	return SetColor(msg, 0, 0, TextCyan)
}

func White(msg string) string {
	return SetColor(msg, 0, 0, TextWhite)
}

func SetColor(msg string, conf, bg, text int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, text, msg, 0x1B)
}
