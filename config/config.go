package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	config Config
	mu     sync.RWMutex
)

type Config struct {
	LogLevel        int
	ServiceLogFile  string
	DatabaseLogFile string

	ServiceName string
	JwtSecret   string

	MongoURL     string
	DatabaseName string
	NumberRetry  int

	ExpiresTimeAccessToken  int
	ExpiresTimeRefreshToken int

	ApiCallBack string
	ApiRedirect string

	AccessTokenType  string
	RefreshTokenType string

	KafkaAddr string
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func Set(c Config) {
	mu.Lock()
	defer mu.Unlock()
	config = c
}

func LoadFromEnv(configPrefix, configSource string) error {
	godotenv.Load(configSource)
	return envconfig.Process(configPrefix, &config)
}
