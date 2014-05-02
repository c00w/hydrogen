package util

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

var defaultdebug *log.Logger = log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

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
	defaultdebug.Output(2, fmt.Sprintf(format, v...))
}

func Short(i string) string {
	i = hex.EncodeToString([]byte(i))
	if len(i) > 8 {
		i = i[:8]
	}
	return i
}
