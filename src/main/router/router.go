package router

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPRouter interface {
	http.Handler
	Add(string, string, http.Handler)
	Set404(http.Handler)
	Set405(http.Handler)
}

type ParamKey struct {
	Name string
}

// New returns router as http.Handler.
func New(ctx context.Context, name string) (HTTPRouter, error) {
	err := fmt.Errorf("%s not implemented", name)
	switch name {
	case "bone":
		return newMuxBone(ctx), nil
	case "httprouter":
		return newMuxHTTPRouter(ctx), nil
	case "vestigo":
		return newMuxVestigo(ctx), nil
	default:
		return nil, err
	}
}
