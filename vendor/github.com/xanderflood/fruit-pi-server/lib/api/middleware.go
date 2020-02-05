package api

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

//DefaultRetryer returns a roundtripper to retry certain response codes
func DefaultRetryer(next RoundTripper) RoundTripper {
	return rehttp.NewTransport(
		next,
		rehttp.RetryAny(
			//Retry all errors (non-responses), TooEarly, TooManyRequests, and 5xx.
			rehttp.RetryTemporaryErr(),
			rehttp.RetryStatuses(http.StatusTooEarly, http.StatusTooManyRequests),
			rehttp.RetryStatusInterval(http.StatusInternalServerError, 999),
		),
		rehttp.ExpJitterDelay(100*time.Millisecond, time.Second),
	)
}

//DefaultAuthorizer returns a roundtripper to attach a bearer token to each request
func DefaultAuthorizer(next RoundTripper, token string) RoundTripper {
	return RoundTripHandler(func(req *http.Request) (*http.Response, error) {
		//TODO copy the body, don't mutate it
		req.Header.Set("Authorization", "Bearer "+token)

		return next.RoundTrip(req)
	})
}
