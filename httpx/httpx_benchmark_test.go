package httpx

import (
	"testing"
)

func BenchmarkRequest_Get(b *testing.B) {
	req := NewClient().R().SetHeaders(map[string]string{
		"header_1":   "value_1",
		"header_2":   "value_2",
		"Connection": "close"}).
		SetHeader("header_1", "value_1").
		SetUserAgent(DefaultUserAgent).
		SetQueryParams(map[string]string{
			"foo1": "bar1",
			"foo2": "bar2",
		})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := req.
				Get("http://localhost:18080/robots.txt")

			if err != nil {
				b.Errorf("Get failed: %v", err)
			}
			//b.Logf("%s -- %d", req.Url(), resp.Status())
		}
	})
}
