package client

import (
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

// NewRedisPool returns connection pool to Redis Server.
func NewRedisPool(addr string, m int, d time.Duration) (*redis.Pool, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &redis.Pool{
			MaxIdle:     m,
			IdleTimeout: d,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", u.Host)
			},
		},
		nil
}

/*
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
