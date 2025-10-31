package logger

import "log"

func Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func Error(msg string) {
	log.Printf("[ERROR] %v\n", msg)
}
