package stress

import (
	"fmt"
	"strings"
	"time"
)

// RequestResult is adapted result from http.Response.
type RequestResult struct {
	statusCode int
	duration   time.Duration
}

// BatchResult is an aggregate for a batch of requests.
type BatchResult struct {
	rr   []RequestResult
	errs []error
}

// NewBatchResult creates new BatchResult with the given state.
func NewBatchResult(rr []RequestResult, errs []error) BatchResult {
	return BatchResult{rr: rr, errs: errs}
}

// NrResponses returns the amount of responses for the batch.
func (br BatchResult) NrResponses() int {
	return len(br.rr)
}

// AverageDuration returns what is the average duration for the given batch.
func (br BatchResult) AverageDuration() time.Duration {
	if len(br.rr) == 0 {
		return 0
	}

	var sum time.Duration
	for _, res := range br.rr {
		sum += res.duration
	}

	return sum / time.Duration(len(br.rr))
}

// ModeStatusCode returns the most frequently occuring status code in the given batch.
func (br BatchResult) ModeStatusCode() int {
	occ := map[int]int{}
	for _, res := range br.rr {
		occ[res.statusCode]++
	}

	var maxK, maxV int
	for k, v := range occ {
		if maxV < v {
			maxV = v
			maxK = k
		}
	}
	return maxK
}

// ErrorRatio returns how many of the HTTP requests failed in the given batch.
func (br BatchResult) ErrorRatio() float64 {
	var nErr float64
	for _, res := range br.rr {
		if res.statusCode < 200 || res.statusCode > 299 {
			nErr++
		}
	}
	return nErr / float64(len(br.rr))
}

// MultiBatchResult is the result for multiple batches of execution.
type MultiBatchResult []BatchResult

// String implements fmt.Stringer.
func (mbr MultiBatchResult) String() string {
	var sb strings.Builder
	for i, br := range mbr {
		sb.WriteString(fmt.Sprintf("Summary for batch %d (%d requests):\n", i+1, len(br.rr)))
		sb.WriteString(fmt.Sprintf("\tMean status code: %d\n", br.ModeStatusCode()))
		sb.WriteString(fmt.Sprintf("\tAverage response time: %s\n", br.AverageDuration()))
		sb.WriteString(fmt.Sprintf("\tError ratio: %f\n", br.ErrorRatio()))
	}
	return sb.String()
}
