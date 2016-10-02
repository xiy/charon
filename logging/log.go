package logging

import (
	"fmt"
	"log"
	"os"

	"comail.io/go/colog"
	gokitlog "github.com/go-kit/kit/log"
)

// DefaultStackDepth is the default level to plumb to when logging
const DefaultStackDepth = 1

// NewStdoutLogger creates a new go-kit based logger
func NewStdoutLogger() (logger gokitlog.Logger) {
	logger = gokitlog.NewLogfmtLogger(os.Stderr)
	logger = gokitlog.NewContext(logger).With("ts", gokitlog.DefaultTimestampUTC)
	logger = gokitlog.NewContext(logger).With("caller", gokitlog.DefaultCaller)
	return
}

// NewCoLogLogger creates a new CoLog based logger
func NewCoLogLogger(domain string) (logger *log.Logger) {
	cl := colog.NewCoLog(os.Stdout, fmt.Sprintf("%s ", domain), log.LstdFlags|log.Lshortfile)

	return cl.NewLogger()
}
