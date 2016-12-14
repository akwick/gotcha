package tests

import (
	"goretech/analysis/worklist"
	"testing"
)

func benchmarkDoAnalysis(file []string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		worklist.DoAnalysis("goretech/analysis", file, "./sourcesAndSinksTest.txt", false, "", true)
	}
}

func BenchmarkStructTest(b *testing.B) {
	benchmarkDoAnalysis([]string{"./exampleCode/structTest.go"}, b)
}
func BenchmarkStructTestV2(b *testing.B) {
	benchmarkDoAnalysis([]string{"./exampleCode/structTestV2.go"}, b)
}
func BenchmarkStructTEstV3Val(b *testing.B) {
	benchmarkDoAnalysis([]string{"./exampleCode/structTestV3Val.go"}, b)
}
func BenchmarkStructTEstV3Ref(b *testing.B) {
	benchmarkDoAnalysis([]string{"./exampleCode/structTestV3Ref.go"}, b)
}
