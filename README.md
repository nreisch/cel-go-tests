# Overview
Experimenting and investigating CEL and the CEL-go engine

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

## Benchmarking Arrays Evaluation
Test Name|Number of Policy Expressions|Array Input Length |Time (ns) per op|Bytes per op|Allocs
---|---|---|---|---|---|
BenchmarkTest_Policy1_Array1_Evaluate-16|1|1|14320 ns/op|8243 B/op|98 allocs/op
BenchmarkTest_Policy100_Array1_Evaluate-16|100|1|1687981 ns/op|829979 B/op|9936 allocs/op
BenchmarkTest_Policy100_Array10_Evaluate-16|100|10|1815824 ns/op|889816 B/op|11267 allocs/op
BenchmarkTest_Policy100_Array100_Evaluate-16|100|100|3079584 ns/op|1436478 B/op|29633 allocs/op
BenchmarkTest_Policy1000_Array1_Evaluate-16|1000|1|15934536 ns/op|8745264 B/op|109681 allocs/op
BenchmarkTest_Policy1000_Array10_Evaluate-16|1000|10|17464611 ns/op|9376829 B/op|123757 allocs/op
BenchmarkTest_Policy1000_Array100_Evaluate-16|1000|100|28511772 ns/op|15166177 B/op|314868 allocs/op
BenchmarkTest_Policy1000_Array1000_Evaluate-16|1000|1000|145696300 ns/op|84891340 B/op|2213605 allocs/op

# References
- https://github.com/google/cel-go/
- https://github.com/google/cel-spec/
- https://github.com/google/cel-policy-templates-go
- https://codelabs.developers.google.com/codelabs/cel-go

# Commands

```console
go test -count=1 -bench=. -benchmem
```