package worker

import (
	"fmt"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/internal/worker/tasks"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/hibiken/asynq"
)

type Worker struct {
	cfg    *config.Config
	logger logging.Logger
}

func New(cfg *config.Config, lg logging.Logger) *Worker {
	return &Worker{
		cfg:    cfg,
		logger: lg,
	}
}

func (w *Worker) InitWorker() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", w.cfg.Redis.Host, w.cfg.Redis.Port), Password: w.cfg.Redis.Password},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("task:deliverer", tasks.HandleExecuteTask)

	if err := srv.Run(mux); err != nil {
		w.logger.Fatalf("could not run server: %v", err)
	}
}
