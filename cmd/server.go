package cmd

import (
	"os"

	"github.com/hamiddarani/anakonda/internal/api"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/db/migration"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/spf13/cobra"
)

type Server struct{}

func (cmd Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.run(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "Run Anakonda server",
		Run:   run,
	}
}

func (cmd *Server) run(cfg *config.Config, trap chan os.Signal) {

	logger := logging.NewLogger(cfg.Logger)

	err := db.InitDb(cfg.Postgres)
	if err != nil {
		logger.Warn(logging.Postgres, logging.Startup, err.Error(), nil)
	} else {
		migration.UP_1()
	}

	err = cache.InitRedis(cfg.Redis)
	if err != nil {
		logger.Warn(logging.Redis, logging.Startup, err.Error(), nil)
	}

	server := api.New(cfg, logger)
	go server.Serve()

	logger.Info(logging.OS, logging.SystemKill, "signal trap", map[logging.ExtraKey]interface{}{"signal": (<-trap).String()})
}
