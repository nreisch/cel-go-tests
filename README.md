# Overview
- Experimenting and investigating CEL and the CEL-go engine

# Benchmarking

## Virtual Machine specs
goos: windows <br>
goarch: amd64 <br>
pkg: github.com/nreisch/cel-go-tests <br>
cpu: Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz

## Benchmarking Hot Path Compilation vs Evaluation
Test Name|Number of Policy Expressions|Time (ns) per op|Bytes per op|Allocs
---|---|---|---|---|
BenchmarkTest_Policy1_CompileEvaluate-16|1|61825 ns/op|23416 B/op|476 allocs/op
BenchmarkTest_Policy1_Evaluate-16|1|10301 ns/op|6486 B/op|60 allocs/op
BenchmarkTest_Policy100_CompileEvaluate-16|100|5825781 ns/op|2343144 B/op|47551 allocs/op
BenchmarkTest_Policy100_Evaluate-16|100|1097703 ns/op|650685 B/op|6062 allocs/op

## Benchmarking Arrays Compilation + Evaluation
Test Name|Number of Policy Expressions|Array Input Length |Time (ns) per op|Bytes per op|Allocs
---|---|---|---|---|---|
BenchmarkTest_Policy1_Array1_CompileEvaluate-16|1|1|117215 ns/op|44185 B/op|906 allocs/op
BenchmarkTest_Policy100_Array1_CompileEvaluate-16|100|1|11662914 ns/op|4420092 B/op|90533 allocs/op
BenchmarkTest_Policy100_Array10_CompileEvaluate-16|100|10|11771656 ns/op|4478050 B/op|91824 allocs/op
BenchmarkTest_Policy100_Array100_CompileEvaluate-16|100|100|13330341 ns/op|5024987 B/op|110140 allocs/op
BenchmarkTest_Policy1000_Array1_CompileEvaluate-16|1000|1|114189911 ns/op|44104144 B/op|904917 allocs/op
BenchmarkTest_Policy1000_Array10_CompileEvaluate-16|1000|10 |116303911 ns/op|44686122 B/op|917940 allocs/op
BenchmarkTest_Policy1000_Array100_CompileEvaluate-16|1000|100|126396388 ns/op|50107981 B/op|1100936 allocs/op
BenchmarkTest_Policy1000_Array1000_CompileEvaluate-16|1000|1000|235579520 ns/op|115632910 B/op|2905112 allocs/op

# References
- https://github.com/google/cel-go/
- https://github.com/google/cel-spec/
- https://github.com/google/cel-policy-templates-go
- https://codelabs.developers.google.com/codelabs/cel-go

# Commands

```console
go test -count=1 -bench=. -benchmem
```