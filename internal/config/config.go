package config

import (
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/logging"
)

type App struct {
	Name string `koanf:"name"`
	Port int    `koanf:"port"`
}

type Controller struct {
	ControllerLeaderRedisKey string `koanf:"controller_leader_redis_key"`
	ControllerPrefixKey      string `koanf:"controller_prefix_key"`
}

type Worker struct{}

type Queue struct {
	Channels map[string]string `konaf:"channels"`
}

type Config struct {
	App        *App            `koanf:"app"`
	Controller *Controller     `koanf:"controller"`
	Worker     *Worker         `koanf:"worker"`
	Queue      *Queue          `koanf:"queue"`
	Logger     *logging.Config `koanf:"logger"`
	Postgres   *db.Config      `koanf:"postgres"`
	Redis      *cache.Config   `koanf:"redis"`
}
