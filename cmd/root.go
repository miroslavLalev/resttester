package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"resttester/stress"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	method       string
	headers      []string
	payload      string
	shouldFollow bool
	maxRedirects int
	insecure     bool
	caCert       string
	dryRun       bool
	timeout      string

	strategy string

	durationPlot string
)

var rootCmd = &cobra.Command{
	Use:   "resttest",
	Short: "Stess test a REST endpoint",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("only endpoint should be provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c := stress.Config{}
		c.Endpoint = args[0]
		c.Method = method
		c.Payload = payload
		c.ShouldFollow = shouldFollow
		c.MaxRedirects = maxRedirects
		c.Insecure = insecure
		c.DryRun = dryRun
		c.Strategy = strategy
		c.DurationPlot = durationPlot

		for _, val := range headers {
			split := strings.Split(val, "=")
			if len(split) != 2 {
				return fmt.Errorf("Invalid header '%s'", val)
			}
			c.Headers = append(c.Headers, split)
		}

		var cert []byte
		if caCert != "" {
			f, err := os.Open(caCert)
			if err != nil {
				return fmt.Errorf("failed to open certificate: %s", err)
			}
			cert, err = ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("failed to read from certificate: %s", err)
			}
		}
		c.CACertificates = cert

		t, err := time.ParseDuration(timeout)
		if err != nil {
			return fmt.Errorf("Invalid timeout format")
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(t))
		defer cancel()

		result, err := stress.NewTester(c).Test(ctxWithInterrupt(ctx))
		if err != nil {
			return err
		}

		_, err = os.Stdout.WriteString(result.String())
		return err
	},
}

func ctxWithInterrupt(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cancel()
	}()
	return ctx
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&method, "request", "X", "GET", "Defines a HTTP method for the endpoint.")
	rootCmd.PersistentFlags().StringArrayVarP(&headers, "header", "H", nil, "Additional headers for the request.")
	rootCmd.PersistentFlags().StringVarP(&payload, "payload", "P", "", "Adds body to the request.")
	rootCmd.PersistentFlags().BoolVarP(&shouldFollow, "location", "L", false, "Whether or not to follow 3xx requests.")
	rootCmd.PersistentFlags().IntVar(&maxRedirects, "max-redirs", 0, "Maximum amount of redirects.")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Proceed if server connection is considered insecure.")
	rootCmd.PersistentFlags().StringVarP(&caCert, "ca-certificates", "c", "", "Specify certificate file for client verification.")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Signifies whether only check request should be made.")
	rootCmd.PersistentFlags().StringVarP(&timeout, "timeout", "t", "", "Specifies the test's running time.")
	rootCmd.PersistentFlags().StringVarP(&durationPlot, "plot", "p", "", "Path to where plot will be generated for the current test.")
	rootCmd.PersistentFlags().StringVarP(&strategy, "strategy", "s", "",
		"Specifies the strategy for determining the amount of parallel requests."+
			"Two types of strategies are currently available - linear and exponential.\n"+
			"(*) Linear strategy is based on the form a*x+b. In order to specify it you should type 'lin[a,b]'.\n"+
			"(*) Exponential strategy is based on the form a*b^x. In order to specify it you should type 'exp[a,b]'.\n")
}

// Execute starts the main command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
