package command

import (
	"main/logger"
	"main/option"
)

type optionReceiver struct {
	log logger.Logger
}

// Logger sets logger from packages on another side.
func Logger(l logger.Logger) option.Fn {
	return func(r option.Receiver) error {
		v, _ := v.(*optionReceiver)
		return v.SetLogger(l)
	}
}
