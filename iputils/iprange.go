package iputils

import (
	"fmt"
	"github.com/projectdiscovery/mapcidr"
	"net/netip"
	"strconv"
	"strings"
)

// IpsFromRange 通过 IP 地址范围生成 IP 地址
// @example: IpsFromRange("192.168.0.100-126")
func IpsFromRange(ipRange string) ([]string, error) {
	if !strings.Contains(ipRange, "-") {
		return nil, fmt.Errorf("invalid ip range: %s", ipRange)
	}

	ips := strings.Split(ipRange, "-")
	if len(ips) != 2 {
		return nil, fmt.Errorf("invalid ip range: %s, len=%d", ipRange, len(ips))
	}

	// 起始IP地址
	_, err := netip.ParseAddr(ips[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ip: %s", ips[0])
	}

	// - 的后面必须也是IP地址，或者IP的最大范围
	_, err = netip.ParseAddr(ips[1])
	if err != nil {
		if maxIp, err := strconv.Atoi(ips[1]); err != nil {
			return nil, fmt.Errorf("invalid ip or not a number: %s", ips[1])
		} else {

			// 将 ips[1] 替换成IP地址格式
			ips0 := strings.Split(ips[0], ".")
			ips0[len(ips0)-1] = strconv.Itoa(maxIp)
			ips[1] = strings.Join(ips0, ".")
		}

	}

	// 生成 IP 地址
	cidrs, err := mapcidr.IpRangeToCIDR(ips[0], ips[1])
	if err != nil || len(cidrs) == 0 {
		return nil, fmt.Errorf("parse cidr error: %s", err)
	}

	var restIps []string
	for _, cidr := range cidrs {
		addrs, _ := mapcidr.IPAddresses(cidr)
		restIps = append(restIps, addrs...)
	}

	return restIps, nil
}
