package logger

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func LoadLogDefinition() {
	debug := false

	flag.BoolVar(&debug, "debug", false, "Enable debugging mode")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})

	if debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	}
	log.SetOutput(os.Stderr)
}
