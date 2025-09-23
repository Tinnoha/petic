package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type configOfRedis struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

// Addr - адрес базы
// Password - пароль
// User - имя пользователя
// DB - индификатор базы
// MaxRetries - максимальное количесвто попвыток подключения
// DialTimeout - время прерыва между попытками подключения
// Timeout - Время для записи и чтения

func NewConfigRedis(addr string, password string, db int, maxretries int, dialTimeout time.Duration, timeout time.Duration) *configOfRedis {
	return &configOfRedis{
		Addr:        addr,
		Password:    password,
		DB:          db,
		MaxRetries:  maxretries,
		DialTimeout: dialTimeout,
		Timeout:     timeout,
	}
}

func NewClientRedis(ctx context.Context, cfg configOfRedis) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Println("fail to connect redis server", err.Error())
		return nil, err
	}

	return db, nil
}
