package main

import "testing"

func BenchmarkTestBasicPolicyTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basic_policy_test()
    }
}