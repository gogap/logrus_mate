package redisio

import (
	"errors"
	"io"

	"github.com/adjust/redis"
	"github.com/adjust/redisio"

	"github.com/gogap/logrus_mate"
)

type RedisIOConfig struct {
	Network  string
	Address  string
	Password string
	Db       int64
	ListName string
}

func init() {
	logrus_mate.RegisterWriter("redisio", NewRedisIOWriter)
}

func NewRedisIOWriter(options *logrus_mate.Options) (writer io.Writer, err error) {
	conf := RedisIOConfig{}

	if options != nil {
		conf.Network = options.GetString("network")
		conf.Address = options.GetString("address")
		conf.Password = options.GetString("password")
		conf.Db = options.GetInt64("db")
		conf.ListName = options.GetString("list-name")
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
