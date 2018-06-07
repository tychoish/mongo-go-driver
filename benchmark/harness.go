package benchmark

import (
	"context"
	"testing"
	"time"
)

const (
	ExecutionTimeout = 5 * time.Minute
	MinimumRuntime   = 20 * time.Second
	MinIterations    = 100

	hundred     = 100
	thousand    = 10 * hundred
	tenThousand = 10 * thousand
)

type BenchCase func(context.Context, int) error
type BenchFunction func(*testing.B)

func WrapCase(bench BenchCase) BenchFunction {
	name := getName(bench)
	return func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		err := bench(ctx, b.N)
		if err != nil {
			b.Fatalf("benchmark %s encountered %s error", name, err.Error())
		}
	}
}

func getAllCases() []*CaseDefinition {
	return []*CaseDefinition{
		{
			Bench:   CanaryIncCase,
			Count:   hundred,
			Size:    -1,
			Runtime: MinimumRuntime,
		},
		{
			Bench:   GlobalCanaryIncCase,
			Count:   hundred,
			Size:    -1,
			Runtime: MinimumRuntime,
		},
	}
}
