package stress

import (
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Plotter is a struct for creating plots for the given results.
type Plotter struct {
	results MultiBatchResult
}

// NewPlotter creates new Plotter for the given results.
func NewPlotter(results MultiBatchResult) *Plotter {
	return &Plotter{results: results}
}

// GenerateDurationPlot creates plot with the average response time over the
// number of requests.
func (p *Plotter) GenerateDurationPlot(location string) error {
	plot, err := plot.New()
	if err != nil {
		return err
	}

	plot.Title.Text = "Response time for number of concurrent requests"
	plot.X.Label.Text = "Number of requests"
	plot.Y.Label.Text = "Response time (ms)"

	pts := make(plotter.XYs, len(p.results))
	for i, res := range p.results {
		pts[i].X = float64(res.NrResponses())
		pts[i].Y = float64(res.AverageDuration() / time.Millisecond)
	}
	plotutil.AddLinePoints(plot, pts)

	return plot.Save(4*vg.Inch, 4*vg.Inch, location)
}
