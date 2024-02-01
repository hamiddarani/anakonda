package cmd

import (
	"os"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/internal/worker"
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/spf13/cobra"
)

type Worker struct{}

func (cmd Worker) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.run(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "worker",
		Short: "Run Anakonda Worker",
		Run:   run,
	}
}

func (cmd *Worker) run(cfg *config.Config, trap chan os.Signal) {

	logger := logging.NewLogger(cfg.Logger)

	err := db.InitDb(cfg.Postgres)
	if err != nil {
		logger.Warn(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	err = cache.InitRedis(cfg.Redis)
	if err != nil {
		logger.Warn(logging.Redis, logging.Startup, err.Error(), nil)
	}

	worker := worker.New(cfg, logger)
	go worker.InitWorker()

	logger.Info(logging.OS, logging.SystemKill, "signal trap", map[logging.ExtraKey]interface{}{"signal": (<-trap).String()})
}
