package server

import (
	"fmt"
	"net/http"
	"time"

	"main/option"
)

var formatErrTypeAssertion = fmt.Sprintf("option receiver type is <%%T>; want <%T>", &server{})

// Address is option for setting that value from packages on another side.
func Address(a string) option.Fn {
	return func(v interface{}) error {
		if s, ok := v.(*server); ok {
			return s.setAddress(a)
		}
		return fmt.Errorf(formatErrTypeAssertion, v)
	}
}

// Timeout is option for setting that value from packages on another side.
func Timeout(t time.Duration) option.Fn {
	return func(v interface{}) error {
		if s, ok := v.(*server); ok {
			return s.setTimeout(t)
		}
		return fmt.Errorf(formatErrTypeAssertion, v)
	}
}

// Handler is option for setting that value from packages on another side.
func Handler(h http.Handler) option.Fn {
	return func(v interface{}) error {
		if s, ok := v.(*server); ok {
			return s.setHandler(h)
		}
		return fmt.Errorf(formatErrTypeAssertion, v)
	}
}

// Logger is option for setting that value from packages on another side.
func Logger(l logger) option.Fn {
	return func(v interface{}) error {
		if s, ok := v.(*server); ok {
			return s.setLogger(l)
		}
		return fmt.Errorf(formatErrTypeAssertion, v)
	}
}
