package initialize

import (
	"TaskHub/user_service/global"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func InitRedis() error{
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     global.Cfg.Redis.Host + ":" + strconv.Itoa(global.Cfg.Redis.Port),
		Password: global.Cfg.Redis.Password,
		DB:       global.Cfg.Redis.DB,
	})
	_, err := global.Redis.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}
