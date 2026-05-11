package config

import (
	"log"
	"log/slog"
)

func init() {
	// init log
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	slog.SetDefault(slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{AddSource: true})))

}
