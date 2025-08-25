package iputils

import "testing"

func TestIpsFromRange(t *testing.T) {
	ips, err := IpsFromRange("192.168.0.10-192.168.0.49")
	if err != nil {
		t.Fatal(err)
	}
	for _, ip := range ips {
		t.Log(ip)
	}

	ips, err = IpsFromRange("192.168.0.100-126")
	if err != nil {
		t.Fatal(err)
	}
	for _, ip := range ips {
		t.Log(ip)
	}
}
