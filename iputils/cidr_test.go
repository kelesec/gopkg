package iputils

import (
	"fmt"
	"testing"
)

func TestIPsFromCIDR(t *testing.T) {
	ips, err := IPsFromCIDR("192.168.0.1/24")
	if err != nil {
		t.Fatal(err)
	}
	for _, ip := range ips {
		t.Log(ip)
	}
}

func TestCIDRFromIps(t *testing.T) {
	ips := []string{
		"192.168.2.68",
		"192.168.2.74",
		"192.168.2.68",
		"192.168.2.888",
		"192.168.2.70",
		"192.168.2.71",
		"192.168.2.70",
		"192.168.2.69",
		"192.168.2.aaa",
		"10.0.0.3",
		"10.0.0.2",
		"192.168.1.1",
		"2.2.2.2",
		"1.1.1.1",
		"172.16.0.1",
		"localhost",
		"https://baidu.com",
	}

	for i := 0; i < 64; i++ {
		ips = append(ips, fmt.Sprintf("8.10.199.%d", i))
	}

	cidrList, err := CIDRFromIps(ips)
	if err != nil {
		t.Fatal(err)
	}

	for cidr, count := range cidrList {
		t.Logf("%s: %d\n", cidr, count)
	}
}
