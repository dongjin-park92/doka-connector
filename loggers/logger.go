package loggers

import (
	"log"
	"os"
)

var GlobalLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
