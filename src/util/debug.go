package util

import (
	"encoding/hex"
	"log"
	"os"
)

var defaultdebug *log.Logger = log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lmicroseconds)

func init() {
	if os.Getenv("LOG") != "DEBUG" {
		fd, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		defaultdebug = log.New(fd, "", 0)
	}
}

func Debugf(format string, v ...interface{}) {
	defaultdebug.Printf(format, v...)
}

func Short(i string) string {
	return hex.EncodeToString([]byte(i))[:8]
}
