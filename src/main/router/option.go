package router

import (
	"main/logger"
	"main/option"
)

type optionReceiver struct {
	muxKind MuxKind
	log     logger.Logger
}

func (o *optionReceiver) SetLogger(l logger.Logger) error {
	o.log = l
	return nil
}

func (o *optionReceiver) SetKind(k MuxKind) error {
	o.muxKind = k
	return nil
}

func (o *optionReceiver) Receive(options ...option.Fn) error {
	return option.Receive(o, options...)
}

func Kind(k MuxKind) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.SetKind(k)
	}
}

func Logger(l logger.Logger) option.Fn {
	return func(r option.Receiver) error {
		v, _ := r.(*optionReceiver)
		return v.SetLogger(l)
	}
}
