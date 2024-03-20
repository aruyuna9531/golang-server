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
}

func Debug(format string, args ...any) {
	DebugL.Output(2, fmt.Sprintf(format, args...))
}

func Info(format string, args ...any) {
	InfoL.Output(2, fmt.Sprintf(format, args...))
}

var Printf = Debug
