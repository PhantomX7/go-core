package redis

import (
	"fmt"
	"log"
	"time"

	redisClient "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v3"
	"github.com/go-redsync/redsync/v3/redis"
	"github.com/go-redsync/redsync/v3/redis/goredis"

	"github.com/PhantomX7/go-core/utility/errors"
)

type Credentials struct {
	Host     string
	Port     string
	Password string
}

type Client interface {
	Get(prefix string, key string) string
	Set(prefix string, key string, value string, expirationTime time.Duration) error
	Delete(prefix string, key string) error
	Ping() error
	Close() error
	NewMutex(key string) *redsync.Mutex
}

type Redis struct {
	client  *redisClient.Client
	redsync *redsync.Redsync
}

func NewClient(credentials Credentials, appEnv string) Client {
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", credentials.Host, credentials.Port),
		Password: credentials.Password,
		DB:       0,
	})
	status := client.Ping()
	if status.Err() != nil {
		if appEnv != "development" {
			log.Panic(status.Err())
		} else {
			log.Println("warning: redis not connected")
		}
	}

	pool := goredis.NewGoredisPool(client)
	rs := redsync.New([]redis.Pool{pool})

	return &Redis{
		client:  client,
		redsync: rs,
	}
}

func (r *Redis) Get(prefix string, key string) string {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return ""
	}

	val, err := r.client.Get(fmt.Sprint(prefix, key)).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func (r *Redis) Set(prefix string, key string, value string, expirationTime time.Duration) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.Set(fmt.Sprint(prefix, key), value, expirationTime).Err()
	if err != nil {
		fmt.Println(err)
		return errors.ErrUnprocessableEntity
	}

	return nil
}

func (r *Redis) Delete(prefix string, key string) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.Del(fmt.Sprint(prefix, key)).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r *Redis) Ping() error {
	pong, err := r.client.Ping().Result()
	if err != nil {
		log.Println("error pinging redis:", err)
		return errors.ErrUnprocessableEntity

	}

	fmt.Println("connected to redis:", pong)
	return nil
}

func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		log.Println("error closing redis:", err)
		return err
	}
	return nil
}

func (r *Redis) NewMutex(key string) *redsync.Mutex {
	return r.redsync.NewMutex(key)
}
