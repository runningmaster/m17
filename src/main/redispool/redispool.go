package redispool

import (
	"context"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

const defaultURL = "redis://localhost:6379"

// New returns connection pool to Redis Server.
func New(options ...func(*redis.Pool) error) (*redis.Pool, error) {
	p := &redis.Pool{}

	err := Address(defaultURL)(p)
	if err != nil {
		return nil, err
	}

	for i := range options {
		err = options[i](p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Address is TCP address to listen on, "redis://localhost:6379" if empty
func Address(a string) func(*redis.Pool) error {
	return func(p *redis.Pool) error {
		u, err := url.Parse(a)
		if err != nil {
			return err
		}
		p.Dial = func() (redis.Conn, error) {
			return redis.Dial("tcp", u.Host)
		}
		return nil
	}
}

// MaxActive sets maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
func MaxActive(m int) func(*redis.Pool) error {
	return func(p *redis.Pool) error {
		p.MaxActive = m
		return nil
	}
}

// MaxIdle sets maximum number of idle connections in the pool.
func MaxIdle(m int) func(*redis.Pool) error {
	return func(p *redis.Pool) error {
		p.MaxIdle = m
		return nil
	}
}

// IdleTimeout closes connections after remaining idle for this duration. If the value
// is zero, then idle connections are not closed. Applications should set
// the timeout to a value less than the server's timeout.
func IdleTimeout(d time.Duration) func(*redis.Pool) error {
	return func(p *redis.Pool) error {
		p.IdleTimeout = d
		return nil
	}
}

// Wait is rule for Get()'s behavior.
// If Wait is true and the pool is at the MaxActive limit, then Get() waits
// for a connection to be returned to the pool before returning.
func Wait(w bool) func(*redis.Pool) error {
	return func(p *redis.Pool) error {
		p.Wait = w
		return nil
	}
}

// contextClientRedisKey is a unique type to prevent assignment.
// Its associated value should be a redis.Pool.
type contextClientRedisKey struct{}

// ContextWithRedisPool returns a new context based on the provided parent ctx.
func ContextWithRedisPool(ctx context.Context, v *redis.Pool) context.Context {
	return context.WithValue(ctx, contextClientRedisKey{}, v)
}

// RedisConnFromContext returns the redis.Conn value associated with the given key.
func RedisConnFromContext(ctx context.Context) redis.Conn {
	if v, ok := ctx.Value(contextClientRedisKey{}).(*redis.Pool); ok {
		return v.Get()
	}
	panic("*redis.Pool not found")
}
