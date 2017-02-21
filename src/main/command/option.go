package command

import (
	"fmt"
	"io/ioutil"

	"main/option"
)

type logger interface {
	Printf(string, ...interface{})
}

type loggerFn func()

func (l loggerFn) Printf(f string, a ...interface{}) {
	fmt.Fprintf(ioutil.Discard, f, a)
}

var log logger = loggerFn(nil)

// Logger sets logger from packages on another side.
func Logger(l logger) option.Fn {
	return func(_ interface{}) error {
		if l == nil {
			return fmt.Errorf("%v logger", l)
		}
		log = l
		return nil
	}
}
