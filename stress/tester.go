package stress

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Tester is a struct for stress testing.
type Tester struct {
	c Config
}

// NewTester creates new stress tester.
func NewTester(c Config) *Tester {
	return &Tester{c: c}
}

// Test performs stress testing for the given config.
func (t *Tester) Test(ctx context.Context) (fmt.Stringer, error) {
	s, err := NewStrategy(t.c.Strategy)
	if err != nil {
		return nil, err
	}
	res := t.beginTest(ctx, s)

	if t.c.DurationPlot != "" {
		pl := NewPlotter(res)
		err := pl.GenerateDurationPlot(t.c.DurationPlot)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (t *Tester) beginTest(ctx context.Context, s *Strategy) MultiBatchResult {
	var result MultiBatchResult
	doneCh := make(chan []error)

	for {
		var mRes sync.Mutex
		resMap := map[*http.Response]time.Duration{}

		go func() {
			errs := s.ExecBatch(func() error {
				client, err := t.prepareClient()
				if err != nil {
					return err
				}
				tc := NewTimedHTTPClient(client)

				req, err := http.NewRequest(t.c.Method, t.c.Endpoint, strings.NewReader(t.c.Payload))
				if err != nil {
					return err
				}
				req = req.WithContext(ctx)

				res, err := tc.Do(req)
				if err != nil {
					return err
				}
				defer res.Body.Close()

				mRes.Lock()
				defer mRes.Unlock()
				resMap[res] = tc.TimeToRespond()
				return nil
			})
			doneCh <- errs
		}()

		select {
		case errs := <-doneCh:
			result = append(result, t.processBatchResult(resMap, errs))
			if t.c.DryRun {
				return result
			}

		case <-ctx.Done():
			return result
		}
	}
}

func (t *Tester) prepareClient() (*http.Client, error) {
	transport := &http.Transport{}

	config := &tls.Config{}
	if t.c.Insecure {
		config.InsecureSkipVerify = t.c.Insecure
		transport.TLSClientConfig = config
	}
	if t.c.CACertificates != nil {
		pool := x509.NewCertPool()
		if pool.AppendCertsFromPEM(t.c.CACertificates) {
			config.RootCAs = pool
		}

		transport.TLSClientConfig = config
	}

	client := &http.Client{
		Transport: transport,
	}
	if !t.c.ShouldFollow {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		if t.c.MaxRedirects > 0 {
			var numRedirects int
			client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
				if numRedirects < t.c.MaxRedirects-1 {
					numRedirects++
					return nil
				}
				return http.ErrUseLastResponse
			}
		}
	}
	return client, nil
}

func (t *Tester) processBatchResult(resMap map[*http.Response]time.Duration, errs []error) BatchResult {
	var rr []RequestResult
	for res, dur := range resMap {
		rr = append(rr, RequestResult{
			statusCode: res.StatusCode,
			duration:   dur,
		})
	}

	return NewBatchResult(rr, errs)
}
