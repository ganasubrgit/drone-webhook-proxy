package proxy

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

func connectToRedis(addr string) (*redis.Client, error) {
	var err error

	for i := 0; i < 10; i++ {
		client := redis.NewClient(&redis.Options{
			Addr: addr,
			Password: "",
			DB: 0,
		})

		_, err = client.Ping().Result()
		if err == nil {
			return client, nil
		}

		logrus.Info("Cannot connect to redis, sleep for 5 seconds before retrying.")
		time.Sleep(time.Second * 5)
	}

	return nil, err
}