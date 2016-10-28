package main

import (
	"charon/config"
	"charon/logging"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: Integrate go-kit request tracing
// TODO: Integrate go-kit service registry/discovery with Consul
// TODO: Integrate go-kit load balancing
// TODO: Escalator routing (slow/fast belt factorio routing)

func main() {
	logger := logging.NewCoLogLogger("charon")

	config, router := configureRouter(logger)

	errorChan := make(chan error)
	shutdownChan, reloadChan := make(chan os.Signal, 1), make(chan os.Signal, 1)

	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	signal.Notify(reloadChan, syscall.SIGUSR1)

	go func() {
		for {
			select {
			case <-shutdownChan:
				errorChan <- fmt.Errorf("Caught INT -> Going away")
			case <-reloadChan:
				logger.Println("info: Caught USR1 -> Reloading config")
				configureRouter(logger)
			}
		}
	}()

	go func() {
		logger.Println("info: I'm ok to go!")
		errorChan <- http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router)
	}()

	logger.Fatalf("error: error=%s", <-errorChan)
}

func configureRouter(logger *log.Logger) (config config.Config, r *Router) {
	config.Load("config.toml", logger)

	timeout, err := time.ParseDuration(config.ServiceTimeout)
	if err != nil {
		logger.Printf("error: %s [%s]", err, config.ServiceTimeout)
		os.Exit(1)
	}

	r = NewRouter(timeout, logger)

	for name, s := range config.Services {
		logger.Printf("[%s]: %s -> %s", name, s.URL, s.Prefix)
		r.AddService(name, s.URL, s.Prefix)
	}

	return
}
