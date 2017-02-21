package server

import (
	"net/http"
	"net/url"
	"time"

	"main/logger"
	"main/option"
)

type optionReceiver struct {
	address string
	timeout time.Duration
	handler http.Handler
	log     logger.Logger
}

func (o *optionReceiver) setAddress(a string) error {
	u, err := url.Parse(a)
	if err != nil {
		return err
	}
	o.address = u.Host
	return nil
}

func (o *optionReceiver) setTimeout(t time.Duration) error {
	o.timeout = t
	return nil
}

func (o *optionReceiver) setHandler(h http.Handler) error {
	o.handler = h
	return nil
}

func (o *optionReceiver) setLogger(l logger.Logger) error {
	o.log = l
	return nil
}

func (o *optionReceiver) Receive(options ...option.Fn) error {
	return option.Receive(o, options...)
}

// Address is option for setting that value from packages on another side.
func Address(a string) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.setAddress(a)
	}
}

// Timeout is option for setting that value from packages on another side.
func Timeout(t time.Duration) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.setTimeout(t)
	}
}

// Handler is option for setting that value from packages on another side.
func Handler(h http.Handler) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.setHandler(h)
	}
}

// Logger is option for setting that value from packages on another side.
func Logger(l logger.Logger) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.setLogger(l)
	}
}
