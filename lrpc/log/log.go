package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level is the log level.
type Level int

// Enums log level constants.
const (
	LevelNil Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String turns the LogLevel to string.
func (lv *Level) String() string {
	return LevelStrings[*lv]
}

// LevelStrings is the map from log level to its string representation.
var LevelStrings = map[Level]string{
	LevelTrace: "trace",
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

// LevelNames is the map from string to log level.
var LevelNames = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

const levelChar = " TDIWEF"

func init() {
	logging.errlevel = LevelError
	logging.outlevel = LevelInfo
}

// Flush flushes all pending log I/O.
func Flush() {
	os.Stdout.Sync()
	os.Stderr.Sync()
	//logging.lockAndFlushAll()
}

// loggingT collects all the global state of the logging setup.
type loggingT struct {
	// mu protects the remaining elements of this structure and is
	// used to synchronize logging.
	mu sync.Mutex
	// file holds writer for each of the log types.

	errlevel Level
	outlevel Level
}

var logging loggingT

func (l *loggingT) setErrLevel(lv Level) {
	l.errlevel = lv
}
func (l *loggingT) setOutLevel(lv Level) {
	l.outlevel = lv
}

func (l *loggingT) getErrLevel() Level {
	return l.errlevel
}
func (l *loggingT) getOutLevel() Level {
	return l.outlevel
}

/*
header formats a log header as defined by the C++ implementation.
It returns a buffer containing the formatted header.

Log lines have this form:

	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...

where the fields are defined as follows:

	L                A single character, representing the log level (eg 'I' for INFO)
	mm               The month (zero padded; ie May is '05')
	dd               The day (zero padded)
	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
	threadid         The space-padded thread ID as returned by GetTID()
	file             The file name
	line             The line number
	msg              The user-supplied message
*/
func (l *loggingT) header(s Level) string {
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	now := time.Now()
	_, file, line, ok := runtime.Caller(3) // It's always the same number of frames to the user's call.
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	if line < 0 {
		line = 0 // not a real line number, but acceptable to someDigits
	}
	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	_, month, day := now.Date()
	hour, minute, second := now.Clock()

	header := fmt.Sprintf("[%c %02d%02d %02d:%02d:%02d %s:%d] ", levelChar[s], int(month), day, hour, minute, second, file, line)
	return header
}

func (l *loggingT) print(s Level, args ...interface{}) {
	header := l.header(s)
	body := fmt.Sprintln(args...)
	l.output(s, header, body)
}

func (l *loggingT) printf(s Level, format string, args ...interface{}) {
	header := l.header(s)
	body := fmt.Sprintf(format+"\n", args...)
	l.output(s, header, body)
}

// output writes the data to the log files and releases the buffer.
func (l *loggingT) output(s Level, header, body string) {
	el, ol := l.getErrLevel(), l.getOutLevel()
	var out *os.File

	if ol <= s && s < el {
		out = os.Stdout
	}
	if s >= el {
		out = os.Stderr
	}

	if out == nil {
		return
	}

	l.mu.Lock()
	fmt.Fprint(out, header, body)

	if s == LevelFatal {
		// Make sure we see the trace for the current goroutine on standard error.
		os.Stderr.Write(stacks(false))
		// Write the stack trace for all goroutines to the files.
		//trace := stacks(true)
		l.mu.Unlock()
		timeoutFlush(10 * time.Second)
		os.Exit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	}
	l.mu.Unlock()
}

// timeoutFlush calls Flush and returns when it completes or after timeout
// elapses, whichever happens first.  This is needed because the hooks invoked
// by Flush may deadlock when glog.Fatal is called from a hook that holds
// a lock.
func timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		Flush() // calls logging.lockAndFlushAll()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
	}
}

// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

func SetErrLevel(s Level) {
	logging.setErrLevel(s)
}
func SetOutLevel(s Level) {
	logging.setOutLevel(s)
}

func Fatal(args ...interface{}) {
	logging.print(LevelFatal, args...)
}
func Fatalf(format string, args ...interface{}) {
	logging.printf(LevelFatal, format, args...)
}

func Error(args ...interface{}) {
	logging.print(LevelError, args...)
}
func Errorf(format string, args ...interface{}) {
	logging.printf(LevelError, format, args...)
}

func Warn(args ...interface{}) {
	logging.print(LevelWarn, args...)
}
func Warnf(format string, args ...interface{}) {
	logging.printf(LevelWarn, format, args...)
}

func Info(args ...interface{}) {
	logging.print(LevelInfo, args...)
}
func Infof(format string, args ...interface{}) {
	logging.printf(LevelInfo, format, args...)
}

func Debug(args ...interface{}) {
	logging.print(LevelDebug, args...)
}
func Debugf(format string, args ...interface{}) {
	logging.printf(LevelDebug, format, args...)
}

func Trace(args ...interface{}) {
	logging.print(LevelTrace, args...)
}
func Tracef(format string, args ...interface{}) {
	logging.printf(LevelTrace, format, args...)
}
