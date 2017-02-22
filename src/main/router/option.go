package router

import (
	"io/ioutil"
	"log"
)

type logger interface {
	Printf(string, ...interface{})
}

type Option struct {
	log  logger
	kind kindMux
}

var defaultOption = &Option{
	log:  log.New(ioutil.Discard, "", 0),
	kind: kindHTTPRouter,
}

func (o *Option) setLogger(l logger) error {
	o.log = l
	return nil
}

func (o *Option) setKind(k kindMux) error {
	o.kind = k
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

func Bone(o *Option) error {
	return o.setKind(kindBone)
}

func HTTPRouter(o *Option) error {
	return o.setKind(kindHTTPRouter)
}

func Vestigo(o *Option) error {
	return o.setKind(kindVestigo)
}
