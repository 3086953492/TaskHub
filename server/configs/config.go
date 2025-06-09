package configs

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	UserDatabase DatabaseConfig `mapstructure:"userDatabase"`
	TaskDatabase DatabaseConfig `mapstructure:"taskDatabase"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Session  SessionConfig  `mapstructure:"session"`
}

var (
	Cfg *Config
	UserDB *gorm.DB
	TaskDB *gorm.DB
)
 
func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	Cfg = &Config{}
	return viper.Unmarshal(Cfg)
}