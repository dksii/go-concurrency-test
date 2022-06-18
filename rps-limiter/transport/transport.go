package transport

import (
	"net/http"

	"golang.org/x/time/rate"
)

type RPSLimitedTransport struct {
	baseRt  http.RoundTripper
	limiter *rate.Limiter
}

func NewRPSLimitedTransport(rt http.RoundTripper, limit uint) *RPSLimitedTransport {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return &RPSLimitedTransport{
		rt,
		rate.NewLimiter(rate.Limit(limit), 1),
	}
}

func (rt *RPSLimitedTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if err := rt.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return rt.baseRt.RoundTrip(req)
}
