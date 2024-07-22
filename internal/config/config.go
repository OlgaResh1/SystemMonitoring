package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf     `mapstructure:"logger"`
	GRPC    GRPCConfig     `mapstructure:"grpc"`
	Stat    StatisticsConf `mapstructure:"statistics"`
	Metrics MetricsConf    `mapstructure:"metrics"`
}
type GRPCConfig struct {
	Address        string        `mapstructure:"addr"`
	RequestTimeout time.Duration `mapstructure:"requesttimeout"`
	MinInterval    time.Duration `mapstructure:"mininterval"`
}
type LoggerConf struct {
	Level     string `mapstructure:"level"`
	Format    string `mapstructure:"format"`
	AddSource bool   `mapstructure:"addsource"`
}
type StatisticsConf struct {
	IntervalStat time.Duration `mapstructure:"intervalstat"`
	MaxAvgWindow time.Duration `mapstructure:"maxavgwindow"`
}
type MetricsConf struct {
	LoadavgEnabled bool `mapstructure:"loadavg"`
	MemoryEnabled  bool `mapstructure:"memory"`
	CPUEnabled     bool `mapstructure:"cpu"`
	DiskEnabled    bool `mapstructure:"disk"`
	NetworkEnabled bool `mapstructure:"network"`
}

func (config *Config) ReadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	return nil
}

func NewConfig() Config {
	return Config{}
}
