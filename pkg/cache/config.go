package cache

import "time"

type Config struct {
	Host         string        `koanf:"host"`
	Password     string        `koanf:"password"`
	Port         int           `koanf:"port"`
	DB           int           `koanf:"db"`
	DialTimeout  time.Duration `koanf:"dial_timeout"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
	PoolSize     int           `koanf:"pool_size"`
	PoolTimeout  time.Duration `koanf:"pool_timeout"`
}
