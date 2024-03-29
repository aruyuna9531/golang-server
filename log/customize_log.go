package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	debugL *log.Logger
	infoL  *log.Logger
	warnL  *log.Logger
	errorL *log.Logger
	fatalL *log.Logger
)

func init() {
	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file: ", err)
	}

	debugL = log.New(os.Stdout, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile)
	infoL = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	warnL = log.New(os.Stdout, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile)
	errorL = log.New(io.MultiWriter(file, os.Stderr), "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
	fatalL = log.New(io.MultiWriter(file, os.Stderr), "[FATAL]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(format string, args ...any) {
	if debugL == nil {
		return
	}
	debugL.Output(2, fmt.Sprintf(format, args...))
}

func Info(format string, args ...any) {
	if infoL == nil {
		return
	}
	infoL.Output(2, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
	if warnL == nil {
		return
	}
	warnL.Output(2, fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
	if errorL == nil {
		return
	}
	errorL.Output(2, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...any) {
	if fatalL == nil {
		return
	}
	fatalInfo := fmt.Sprintf(format, args...)
	fatalL.Output(2, fatalInfo)
	panic(fatalInfo)
}
