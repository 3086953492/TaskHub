package database

import (
	"fmt"
	"TaskHub/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Viper.GetString("db.username"),
		config.Viper.GetString("db.password"),
		config.Viper.GetString("db.addr"),
		config.Viper.GetString("db.name"))

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
