package api

import (
	"io/ioutil"
	"log"

	"github.com/garyburd/redigo/redis"
)

/*
	api.Router(r),
	api.Redis(p),
*/

type logger interface {
	Printf(string, ...interface{})
}

type Option struct {
	log       logger
	redisPool *redis.Pool
}

var defaultOption = &Option{
	log: log.New(ioutil.Discard, "", 0),
}

func (o *Option) setLogger(l logger) error {
	o.log = l
	return nil
}

func (o *Option) setRedis(p *redis.Pool) error {
	o.redisPool = p
	return nil
}

func (o *Option) override(options ...func(*Option) error) error {
	var err error
	for i := range options {
		err = options[i](o)
		if err != nil {
			return err
		}
	}
	return nil
}

func Logger(l logger) func(*Option) error {
	return func(o *Option) error {
		return o.setLogger(l)
	}
}

func Redis(p *redis.Pool) func(*Option) error {
	return func(o *Option) error {
		return o.setRedis(p)
	}
}
