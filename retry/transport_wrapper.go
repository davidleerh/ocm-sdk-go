/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the implementations of a transport wrapper that knows how
// to retry requests.

package retry

import (
	"context"
	"math/rand"
	"strings"

	"fmt"
	"net/http"
	"time"

	"github.com/openshift-online/ocm-sdk-go/logging"
)

// Default configuration:
const (
	DefaultLimit    = 2
	DefaultInterval = 1 * time.Second
	DefaultJitter   = 0.2
)

// TransportWrapperBuilder contains the data and logic needed to create a new retry transport
// wrapper.
type TransportWrapperBuilder struct {
	logger   logging.Logger
	limit    int
	interval time.Duration
	jitter   float64
}

// TransportWrapper contains the data and logic needed to wrap an HTTP round tripper with another
// one that adds retry capability.
type TransportWrapper struct {
	logger   logging.Logger
	limit    int
	interval time.Duration
	jitter   float64
}

// roundTripper is a round tripper that adds retry logic.
type roundTripper struct {
	logger    logging.Logger
	limit     int
	interval  time.Duration
	jitter    float64
	transport http.RoundTripper
}

// Make sure that we implement the interface:
var _ http.RoundTripper = (*roundTripper)(nil)

// NewTransportWrapper creates a new builder that can then be used to configure and create a new
// retry round tripper.
func NewTransportWrapper() *TransportWrapperBuilder {
	return &TransportWrapperBuilder{
		limit:    DefaultLimit,
		interval: DefaultInterval,
		jitter:   DefaultJitter,
	}
}

// Logger sets the logger that will be used by the wrapper and by the round trippers that it
// creates.
func (b *TransportWrapperBuilder) Logger(value logging.Logger) *TransportWrapperBuilder {
	b.logger = value
	return b
}

// Limit sets the maximum number of retries for a request. When this is zero no retries will be
// performed. The default value is two.
func (b *TransportWrapperBuilder) Limit(value int) *TransportWrapperBuilder {
	b.limit = value
	return b
}

// Interval sets the time to wait before the first retry. The interval time will be doubled for each
// retry. For example, if this is set to one second then the first retry will happen approximately
// one second after the failure of the initial request, the second retry will happen affer four
// seconds, the third will happen after eitght seconds, so on.
func (b *TransportWrapperBuilder) Interval(value time.Duration) *TransportWrapperBuilder {
	b.interval = value
	return b
}

// Jitter sets a factor that will be used to randomize the retry intervals. For example, if this is
// set to 0.1 then a random adjustment between -10% and +10% will be done to the interval for each
// retry.  This is intended to reduce simultaneous retries by clients when a server starts failing.
// The default value is 0.2.
func (b *TransportWrapperBuilder) Jitter(value float64) *TransportWrapperBuilder {
	b.jitter = value
	return b
}

// Build uses the information stored in the builder to create a new transport wrapper.
func (b *TransportWrapperBuilder) Build(ctx context.Context) (result *TransportWrapper, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("logger is mandatory")
		return
	}
	if b.limit < 0 {
		err = fmt.Errorf(
			"retry limit %d isn't valid, it should be greater or equal than zero",
			b.limit,
		)
		return
	}
	if b.interval <= 0 {
		err = fmt.Errorf(
			"retry interval %s isn't valid, it should be greater than zero",
			b.interval,
		)
		return
	}
	if b.jitter < 0 || b.jitter > 1 {
		err = fmt.Errorf(
			"retry jitter %f isn't valid, it should be between zero and one",
			b.jitter,
		)
		return
	}

	// Create and populate the object:
	result = &TransportWrapper{
		logger:   b.logger,
		limit:    b.limit,
		interval: b.interval,
		jitter:   b.jitter,
	}

	return
}

// Wrap creates a new round tripper that wraps the given one and implements the retry logic.
func (w *TransportWrapper) Wrap(transport http.RoundTripper) http.RoundTripper {
	return &roundTripper{
		logger:    w.logger,
		limit:     w.limit,
		interval:  w.interval,
		jitter:    w.jitter,
		transport: transport,
	}
}

// Limit returns the maximum number of retries.
func (w *TransportWrapper) Limit() int {
	return w.limit
}

// Interval returns the initial retry interval.
func (w *TransportWrapper) Interval() time.Duration {
	return w.interval
}

// Jitter returns the retry interval jitter factor.
func (w *TransportWrapper) Jitter() float64 {
	return w.jitter
}

// Close releases all the resources used by the wrapper.
func (w *TransportWrapper) Close() error {
	return nil
}

// RoundTrip is the implementation of the round tripper interface.
func (t *roundTripper) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// Get the context:
	ctx := request.Context()

	// Calculate the max number of times that we can retry this request:
	limit := 0
	if t.retriable(request) && t.limit > limit {
		limit = t.limit
	}

	// Try to send the request till it succeeds or else the retry limit is exceeded:
	attempt := 0
	for {
		// If this is not the first attempt then we should wait:
		if attempt > 0 {
			t.sleep(ctx, attempt)
		}

		// Do an attempt, and return inmediately if this is the last one:
		response, err = t.transport.RoundTrip(request)
		attempt++
		if attempt > limit {
			return
		}

		// Handle errors without HTTP response:
		if err != nil {
			message := err.Error()
			switch {
			case strings.Contains(message, "PROTOCOL_ERROR"):
				// If the request failed due to a protocol error then we should retry:
				t.logger.Warn(
					ctx,
					"Request for method %s and URL '%s' failed with protocol error, "+
						"will try again: %v",
					request.Method, request.URL, err,
				)
				continue
			default:
				// For any other error we just report it to the caller:
				err = fmt.Errorf("can't send request: %w", err)
				return
			}
		}

		// Handle HTTP responses with error codes:
		code := response.StatusCode
		switch {
		case code == http.StatusServiceUnavailable || code == http.StatusTooManyRequests:
			// For 429 and 503 we know that the server didn't process the request, so we
			// can safely retry.
			t.logger.Warn(
				ctx,
				"Request for method %s and URL '%s' failed with code %d, "+
					"will try again",
				request.Method, request.URL, code,
			)
			continue
		default:
			// For any other status code we can't be sure if the server processed the
			// request, so we just return the result to the caller.
			return
		}
	}
}

// retrieable decides if retrying is acceptable for the given request.
func (t *roundTripper) retriable(request *http.Request) bool {
	if request == nil {
		return false
	}
	methodGood := request.Method == http.MethodGet || request.Method == http.MethodDelete
	bodyGood := request.Body == nil
	return methodGood && bodyGood
}

// sleep calculates a retry interval taking into account the configured interval and jitter factor
// and then waits that time.
func (t *roundTripper) sleep(ctx context.Context, attempt int) {
	// Start with the configured interval:
	interval := t.interval

	// Double the interval for each attempt:
	interval *= 1 << attempt

	// Adjust the interval adding or subtracting a random amount. For example, if the jitter
	// factor given in the configuration is 0.1 will add or sustract up to a 10%.
	factor := t.jitter * (1 - 2*rand.Float64())
	delta := time.Duration(float64(interval) * factor)
	interval += delta

	// Go sleep for a while:
	t.logger.Debug(ctx, "Wating %s before next attempt", interval)
	time.Sleep(interval)
}
