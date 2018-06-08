package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	ExecutionTimeout = 5 * time.Minute
	StandardRuntime  = time.Minute
	MinimumRuntime   = 10 * time.Second
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
		require.NoError(b, err, "case='%s'", name)
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
		{
			Bench:   BSONFlatDocumentEncoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONFlatDocumentDecodingLazy,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONFlatDocumentDecoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONDeepDocumentEncoding,
			Count:   tenThousand,
			Size:    19640000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONDeepDocumentDecodingLazy,
			Count:   tenThousand,
			Size:    19640000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONDeepDocumentDecoding,
			Count:   tenThousand,
			Size:    19640000,
			Runtime: StandardRuntime,
		},
		// {
		//	Bench:   BSONFullDocumentEncoding,
		//	Count:   tenThousand,
		//	Size:    57340000,
		//	Runtime: StandardRuntime,
		// },
		// {
		//	Bench:   BSONFullDocumentDecodingLazy,
		//	Count:   tenThousand,
		//	Size:    57340000,
		//	Runtime: StandardRuntime,
		// },
		// {
		//	Bench:   BSONFullDocumentDecoding,
		//	Count:   tenThousand,
		//	Size:    57340000,
		//	Runtime: StandardRuntime,
		// },
		{
			Bench:   BSONFlatReaderDecoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONDeepReaderDecoding,
			Count:   tenThousand,
			Size:    19640000,
			Runtime: StandardRuntime,
		},
		// {
		//	Bench:   BSONFullReaderDecoding,
		//	Count:   tenThousand,
		//	Size:    57340000,
		//	Runtime: StandardRuntime,
		// },
		{
			Bench:   BSONFlatMapDecoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONDeepMapDecoding,
			Count:   tenThousand,
			Size:    19640000,
			Runtime: StandardRuntime,
		},
		// {
		//	Bench:   BSONFullMapDecoding,
		//	Count:   tenThousand,
		//	Size:    57340000,
		//	Runtime: StandardRuntime,
		// },
		{
			Bench:   BSONFlatStructDecoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONFlatStructTagsDecoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONFlatStructEncoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
		{
			Bench:   BSONFlatStructTagsEncoding,
			Count:   tenThousand,
			Size:    75310000,
			Runtime: StandardRuntime,
		},
	}
}
