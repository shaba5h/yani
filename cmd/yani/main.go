package main

import (
	"github.com/shaba5h/yani/internal/config"
	"github.com/shaba5h/yani/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("Starting app", "env", cfg.Env)
}
