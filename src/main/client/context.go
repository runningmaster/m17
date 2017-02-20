package client

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

// contextClientRedisKey is a unique type to prevent assignment.
// Its associated value should be a redis.Pool.
type contextClientRedisKey struct{}

// WithRedisPool returns a new context based on the provided parent ctx.
func WithRedisPool(ctx context.Context, v *redis.Pool) context.Context {
	ctxKey := contextClientRedisKey{}
	return context.WithValue(ctx, ctxKey, v)
}

// ContextRedisConn returns the redis.Conn value associated with the given key.
func ContextRedisConn(ctx context.Context) redis.Conn {
	ctxKey := contextClientRedisKey{}
	if v, ok := ctx.Value(ctxKey).(*redis.Pool); ok {
		return v.Get()
	}
	panic("redis.Pool not found")
}
