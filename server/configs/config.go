package configs

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Config struct {
	Server       ServerConfig   `mapstructure:"server"`
	UserDatabase DatabaseConfig `mapstructure:"userDatabase"`
	TaskDatabase DatabaseConfig `mapstructure:"taskDatabase"`
	Redis        RedisConfig    `mapstructure:"redis"`
	JWT          JWTConfig      `mapstructure:"jwt"`
	Log          LogConfig      `mapstructure:"log"`
	Session      SessionConfig  `mapstructure:"session"`
}

var (
	Cfg    *Config
	UserDB *gorm.DB
	TaskDB *gorm.DB
	Validate *validator.Validate
)

func InitConfig() error {
	viper.AddConfigPath("./configs")

	// 优先读取YAML配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {

		// 如果YAML读取失败，尝试读取JSON配置
		fmt.Printf("YAML配置读取失败，尝试读取JSON配置: %v\n", err)

		viper.SetConfigName("config")
		viper.SetConfigType("json")

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("配置文件读取失败，YAML和JSON配置都无法读取: %v", err)
		}

		fmt.Println("成功读取JSON配置文件")
	} else {
		fmt.Println("成功读取YAML配置文件")
	}

	Cfg = &Config{}
	return viper.Unmarshal(Cfg)
}
