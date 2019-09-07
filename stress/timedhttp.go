package stress

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"time"
)

// TimedHTTPClient is decorator around http.Client that traces time to respond.
type TimedHTTPClient struct {
	*http.Client

	start time.Time
	ttb   time.Time
	dnss  time.Time
	dnsd  time.Time
}

// NewTimedHTTPClient creates client from the origin one.
func NewTimedHTTPClient(client *http.Client) *TimedHTTPClient {
	return &TimedHTTPClient{Client: client}
}

// TimeToRespond returns how long it took the client to respond.
//
// DNS lookup is removed from that time because there are caches which
// could offset the average time.
func (t *TimedHTTPClient) TimeToRespond() time.Duration {
	// remove the dns lookup time
	return t.ttb.Sub(t.start) - t.dnsd.Sub(t.dnss)
}

// Do implements http.Client
func (t *TimedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.withTimedContext(req.Context()))

	t.start = time.Now()
	return t.Client.Do(req)
}

func (t *TimedHTTPClient) withTimedContext(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GetConn: func(_ string) {
			t.reset() // always reset on first call
			t.start = time.Now()
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t.dnss = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t.dnsd = time.Now()
		},
		GotFirstResponseByte: func() {
			t.ttb = time.Now()
		},
	})
}

func (t *TimedHTTPClient) reset() {
	t.start = time.Time{}
	t.ttb = time.Time{}
	t.dnss = time.Time{}
	t.dnsd = time.Time{}
}
