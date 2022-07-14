package fetch

import (
	"net/http"
)

type Option func(o *opts)

type opts struct {
	httpClient *http.Client
}

func defaultOpts() opts {
	return opts{
		httpClient: http.DefaultClient,
	}
}

func resolveOpts(options ...Option) opts {
	o := defaultOpts()
	for _, option := range options {
		option(&o)
	}
	return o
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *opts) {
		o.httpClient = client
	}
}
