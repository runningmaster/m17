package redis

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

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
	panic("redis.Pool not found")
}
