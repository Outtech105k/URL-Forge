package utils

import "time"

type Config struct {
	ServerEndpoint  string        `env:"ENDPOINT"`
	ServerPort      string        `env:"PORT" envDefault:"8080"`
	RedisAddr       string        `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisPassword   string        `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB         int           `env:"REDIS_DB" envDefault:"0"`
	AppName         string        `env:"APP_NAME" envDefault:"URL Forge"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	OGPFetchTimeout time.Duration `env:"OGP_FETCH_TIMEOUT" envDefault:"5s"`
	DefaultIDLength uint32        `env:"DEFAULT_ID_LENGTH" envDefault:"6"`
	MaxIDLength     uint32        `env:"MAX_ID_LENGTH" envDefault:"100"`
	AllowOrigins    string        `env:"ALLOW_ORIGINS" envDefault:"*"`
	MaxRetryCount   int           `env:"MAX_RETRY_COUNT" envDefault:"10"`
	BotUserAgents   []string      `env:"BOT_USER_AGENTS" envSeparator:"," envDefault:"bot,crawler,spider,facebookexternalhit,twitterbot,slackbot,discordbot,whatsapp,line-poker"`
}
