package iputils

import (
	"net"
	"testing"
)

func TestFormat(t *testing.T) {
	values := NewFormat([]string{
		"192.168.2.68",
		"192.168.2.74",
		"192.168.2.68",
		"192.168.2.888",
		"192.168.2.70",
		"192.168.2.71",
		"192.168.2.70",
		"192.168.2.69",
		"192.168.2.aaa",
		"10.0.0.1",
		"192.168.1.1",
		"2.2.2.2",
		"1.1.1.1",
		"172.16.0.1",
		"localhost",
		"https://baidu.com",
	}).Verify(nil).
		Deduplicate().
		Sort().
		Values()

	for _, v := range values {
		t.Log(v)
	}

	t.Log("Private ips:")
	privateIps := NewFormat(values).FilterPrivateIPs()
	for _, v := range privateIps {
		t.Log(v)
	}

	t.Log("Public ips:")
	publicIps := NewFormat(values).FilterValues(func(value string) bool {
		if ip := net.ParseIP(value); ip != nil && !ip.IsPrivate() {
			return true
		}
		return false
	})
	for _, v := range publicIps {
		t.Log(v)
	}
}
