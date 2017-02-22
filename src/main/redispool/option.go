package redispool

import (
	"io/ioutil"
	"log"
	"time"
)

type logger interface {
	Printf(string, ...interface{})
}

type Option struct {
	log     logger
	address string
	maxIdle int
	timeout time.Duration
}

var defaultOption = &Option{
	log:     log.New(ioutil.Discard, "", 0),
	address: "redis://127.0.0.1:6379",
	maxIdle: 128,
	timeout: 60 * time.Second,
}

func (o *Option) setLogger(l logger) error {
	o.log = l
	return nil
}

func (o *Option) setAddress(a string) error {
	o.address = a
	return nil
}

func (o *Option) setMaxIdle(m int) error {
	o.maxIdle = m
	return nil
}

func (o *Option) setTimeout(t time.Duration) error {
	o.timeout = t
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

func Address(a string) func(*Option) error {
	return func(o *Option) error {
		return o.setAddress(a)
	}
}

func MaxIdle(m int) func(*Option) error {
	return func(o *Option) error {
		return o.setMaxIdle(m)
	}
}

func Timeout(t time.Duration) func(*Option) error {
	return func(o *Option) error {
		return o.setTimeout(t)
	}
}
