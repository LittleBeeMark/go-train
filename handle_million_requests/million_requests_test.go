package handle_million_requests

import "testing"

func BenchmarkDealW1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DealW1(200)
	}
}

func BenchmarkDealW2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DealW2(200)
	}
}
