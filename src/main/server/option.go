package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type logger interface {
	Printf(string, ...interface{})
}

type Option struct {
	log     logger
	address string
	timeout time.Duration
	handler http.Handler
}

var defaultOption = &Option{
	log:     log.New(ioutil.Discard, "", 0),
	address: "http://127.0.0.1:8080",
	timeout: 60 * time.Second,
	handler: http.DefaultServeMux,
}

func (o *Option) setLogger(l logger) error {
	o.log = l
	return nil
}

func (o *Option) setAddress(a string) error {
	o.address = a
	return nil
}

func (o *Option) setTimeout(t time.Duration) error {
	o.timeout = t
	return nil
}

func (o *Option) setHandler(h http.Handler) error {
	o.handler = h
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

func Timeout(t time.Duration) func(*Option) error {
	return func(o *Option) error {
		return o.setTimeout(t)
	}
}

func Handler(h http.Handler) func(*Option) error {
	return func(o *Option) error {
		return o.setHandler(h)
	}
}
