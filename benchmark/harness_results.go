package benchmark

import (
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
)

type BenchResult struct {
	Name       string
	Trials     int
	Duration   time.Duration
	Raw        []Result
	DataSize   int
	Operations int
	hasErrors  *bool
}

func (r *BenchResult) EvergreenPerfFormat() (interface{}, error) {
	timings := r.timings()

	median, err := stats.Median(timings)
	if err != nil {
		return nil, err
	}

	min, err := stats.Min(timings)
	if err != nil {
		return nil, err
	}

	max, err := stats.Max(timings)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name": r.Name,
		"results": map[string]interface{}{
			"1": map[string]interface{}{
				"trials":         len(timings),
				"ops_per_trial":  r.Operations,
				"seconds":        r.Duration.Round(time.Millisecond).Seconds(),
				"ops_per_second": r.getThroughput(median),
				"ops_per_second_values": []float64{
					r.getThroughput(min),
					r.getThroughput(max),
				},
			},
		},
	}, nil
}

func (r *BenchResult) timings() []float64 {
	out := []float64{}
	for _, r := range r.Raw {
		out = append(out, r.Duration.Seconds())
	}
	return out
}

func (r *BenchResult) getThroughput(data float64) float64 {
	if r.DataSize > 0 {
		return float64(r.DataSize) / data
	}

	return float64(r.Operations) / data
}

func (r *BenchResult) String() string {
	return fmt.Sprintf("name=%s, trials=%d, secs=%s", r.Name, r.Trials, r.Duration)
}

func (r *BenchResult) HasErrors() bool {
	if r.hasErrors == nil {
		var val bool
		for _, res := range r.Raw {
			if res.Error != nil {
				val = true
				break
			}
		}
		r.hasErrors = &val
	}

	return *r.hasErrors
}

type Result struct {
	Duration   time.Duration
	Iterations int
	Error      error
}
