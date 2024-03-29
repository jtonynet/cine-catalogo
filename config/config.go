package config

import "github.com/spf13/viper"

type API struct {
	Name          string `mapstructure:"API_NAME"`
	Port          string `mapstructure:"PORT"`
	TagVersion    string `mapstructure:"API_TAG_VERSION"`
	Env           string `mapstructure:"API_ENV"`
	Host          string `mapstructure:"HOST"`
	StaticsDir    string `mapstructure:"API_STATICS_DIR"`
	PostersDir    string `mapstructure:"API_POSTERS_DIR"`
	MetricEnabled bool   `mapstructure:"API_METRICS_ENABLED"`
	LogLevel      string `mapstructure:"API_LOG_LEVEL"`
}

type Database struct {
	Driver string `mapstructure:"DATABASE_DRIVER"`

	Port string `mapstructure:"DATABASE_PORT"`
	Host string `mapstructure:"DATABASE_HOST"`
	User string `mapstructure:"DATABASE_USER"`
	Pass string `mapstructure:"DATABASE_PASSWORD"`
	DB   string `mapstructure:"DATABASE_DB"`

	MetricEnabled         bool   `mapstructure:"DATABASE_METRICS_ENABLED"`
	MetricDBName          string `mapstructure:"DATABASE_METRICS_NAME"`
	MetricRefreshInterval uint32 `mapstructure:"DATABASE_METRICS_INTERVAL_IN_SEC"`
	MetricStartServer     bool   `mapstructure:"DATABASE_METRICS_START_SERVER"`
	MetricServerPort      uint32 `mapstructure:"DATABASE_METRICS_SERVER_PORT"`
}

type Config struct {
	API      API      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
