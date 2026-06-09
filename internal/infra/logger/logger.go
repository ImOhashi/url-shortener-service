package logger

import "log"

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func Info(format string, args ...any) {
	log.Printf(colorGreen+"[INFO]"+colorReset+" "+format, args...)
}

func Warning(format string, args ...any) {
	log.Printf(colorYellow+"[WARNING]"+colorReset+" "+format, args...)
}

func Error(format string, args ...any) {
	log.Printf(colorRed+"[ERROR]"+colorReset+" "+format, args...)
}
