package logging

import (
	"os"

	"github.com/go-kit/kit/log"
)

func NewStdoutLogger() (logger log.Logger) {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
	logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	return
}
