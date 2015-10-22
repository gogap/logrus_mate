package redisio

import (
	"errors"
	"io"

	"github.com/adjust/redis"
	"github.com/adjust/redisio"

	"github.com/gogap/logrus_mate"
)

type RedisIOConfig struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Db       int64  `json:"db"`
	ListName string `json:"list_name"`
}

func init() {
	logrus_mate.RegisterWriter("redisio", NewRedisIOWriter)
}

func NewRedisIOWriter(options logrus_mate.Options) (writer io.Writer, err error) {
	conf := RedisIOConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	if conf.ListName == "" {
		err = errors.New("logurs mate: redisio's list name is empty")
		return
	}

	if conf.Network == "" {
		conf.Network = "tcp"
	}

	if conf.Address == "" {
		conf.Address = "127.0.0.1:6379"
	}

	var redisCli *redis.Client

	redisOpt := &redis.Options{
		Network:  conf.Network,
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.Db,
	}

	redisCli = redis.NewClient(redisOpt)

	writer, err = redisio.NewWriter(redisCli, conf.ListName)
	return
}
