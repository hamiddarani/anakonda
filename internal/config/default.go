package config

import (
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/logging"
)

func Default() *Config {
	return &Config{
		App: &App{
			Name: "anakonda",
			Port: 8080,
		},
		Logger: &logging.Config{
			Logger:      "zap",
			Development: true,
			Encoding:    "console",
			Level:       "debug",
		},
		Postgres: &db.Config{
			Host:            "localhost",
			Port:            "5432",
			User:            "postgres",
			Password:        "admin",
			SSLMode:         "disable",
			DbName:          "test",
			MaxIdleConns:    15,
			ConnMaxLifetime: 5,
			MaxOpenConns:    100,
		},
		Redis: &cache.Config{
			Host:         "localhost",
			Port:         6379,
			Password:     "secret",
			DialTimeout:  5,
			ReadTimeout:  5,
			WriteTimeout: 5,
			DB:           0,
			PoolSize:     10,
			PoolTimeout:  15,
		},
		Controller: &Controller{
			ControllerLeaderRedisKey: "anakonda_controller_leader",
			ControllerPrefixKey:      "anakonda_controller_",
		},
		Queue: &Queue{
			Channels: map[string]string{
				"anakonda_new_task": "anakonda_new_task",
			},
		},
	}
}
