package stress

// Config defines configurations for the performed stress test.
type Config struct {
	Endpoint string
	Method   string
	Headers  [][]string
	Payload  string

	ShouldFollow bool
	MaxRedirects int

	CACertificates []byte
	Insecure       bool

	DryRun          bool
	FailureTreshold float32

	Strategy string

	DurationPlot string
}
