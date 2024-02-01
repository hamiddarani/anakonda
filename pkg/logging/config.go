package logging

type Config struct {
	Logger      string `koanf:"zap"`
	Development bool   `koanf:"development"`
	Encoding    string `koanf:"encoding"`
	Level       string `koanf:"level"`
}
