package redispool

import (
	"context"
	"net/url"

	"github.com/garyburd/redigo/redis"
)

// New returns connection pool to Redis Server.
func New(_ context.Context, options ...func(*Option) error) (*redis.Pool, error) {
	err := defaultOption.override(options...)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(defaultOption.address)
	if err != nil {
		return nil, err
	}

	return &redis.Pool{
			MaxIdle:     defaultOption.maxIdle,
			IdleTimeout: defaultOption.timeout,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", u.Host)
			},
		},
		nil
}

/*
// TODO: change *log.Logger to interface when Go1.9 will released.
func WaitRedisFromDisk(p *redis.Pool, d time.Duration, l *log.Logger) error {
	c := p.Get()
	defer c.Close()

	t := time.NewTicker(d)
	defer t.Stop()

	var err error
	for range t.C {
		_, err = c.Do("PING")
		if err != nil {
			if l != nil {
				l.Println(err)
			}
			continue
		}
		break
	}
	return nil
}
*/
