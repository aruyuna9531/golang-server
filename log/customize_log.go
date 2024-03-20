package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	DebugL *log.Logger
	InfoL  *log.Logger
	WarnL  *log.Logger
	ErrorL *log.Logger
	FatalL *log.Logger
)

func init() {
	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file: ", err)
	}

	DebugL = log.New(os.Stdout, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile)
	InfoL = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	WarnL = log.New(os.Stdout, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorL = log.New(io.MultiWriter(file, os.Stderr), "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
	FatalL = log.New(io.MultiWriter(file, os.Stderr), "[FATAL]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(format string, args ...any) {
	if DebugL == nil {
		return
	}
	DebugL.Output(2, fmt.Sprintf(format, args...))
}

func Info(format string, args ...any) {
	if InfoL == nil {
		return
	}
	InfoL.Output(2, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
	if WarnL == nil {
		return
	}
	WarnL.Output(2, fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
	if ErrorL == nil {
		return
	}
	ErrorL.Output(2, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...any) {
	if FatalL == nil {
		return
	}
	fatalInfo := fmt.Sprintf(format, args...)
	FatalL.Output(2, fatalInfo)
	panic(fatalInfo)
}
