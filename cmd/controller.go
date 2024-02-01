package cmd

import (
	"os"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/internal/controller"
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/spf13/cobra"
)

type Controller struct{}

func (cmd Controller) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.run(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "controller",
		Short: "Run Anakonda Controller",
		Run:   run,
	}
}

func (cmd *Controller) run(cfg *config.Config, trap chan os.Signal) {

	logger := logging.NewLogger(cfg.Logger)

	err := db.InitDb(cfg.Postgres)
	if err != nil {
		logger.Warn(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	err = cache.InitRedis(cfg.Redis)
	if err != nil {
		logger.Warn(logging.Redis, logging.Startup, err.Error(), nil)
	}

	controller := controller.New(cfg, logger)
	go controller.InitController()

	logger.Info(logging.OS, logging.SystemKill, "signal trap", map[logging.ExtraKey]interface{}{"signal": (<-trap).String()})
}
