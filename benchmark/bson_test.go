package benchmark

import "testing"

func BenchmarkFlatBSONEncodingDocument(b *testing.B) { WrapCase(FlatBSONEncodingDocument)(b) }
