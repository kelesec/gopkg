package iputils

import (
	"fmt"
	"github.com/projectdiscovery/mapcidr"
	"net"
	"sort"
	"strconv"
	"strings"
)

// IPsFromCIDR 通过 CIDR 生成 IP 地址
// @example: IPsFromCIDR("192.168.1.1/24")
func IPsFromCIDR(cidr string) ([]string, error) {
	// 先解析是否是 CIDR 格式
	if _, _, err := net.ParseCIDR(cidr); err != nil {
		return nil, fmt.Errorf("invalid cidr: %s", cidr)
	}

	return mapcidr.IPAddresses(cidr)
}

// CIDRFromIps 通过 IP_V4 地址生成 CIDR
// return: 返回一个 map，其中 key 是对应的 CIDR，value 是此 CIDR 包含的 IP 个数
func CIDRFromIps(ips []string) (map[string]int, error) {
	// 提取出 IP_V4 地址列表
	ip4List := NewFormat(ips).Deduplicate().Verify(nil).Sort().FilterValues(func(value string) bool {
		return net.ParseIP(value).To4() != nil
	})

	// 将同C段的IP分到一组
	ip4Map := make(map[string][]int)
	for _, ip4 := range ip4List {
		splits := strings.Split(ip4, ".")
		key := strings.Join(splits[:3], ".")
		suffix, _ := strconv.Atoi(splits[3])
		ip4Map[key] = append(ip4Map[key], suffix)
	}

	// 对每个C段进行处理
	cidrMap := make(map[string]int)
	for key, suffixes := range ip4Map {
		sort.Ints(suffixes)

		// 将连续的IP段转CIDR
		startIndex := 0
		for i := 0; i < len(suffixes); i++ {
			if i == len(suffixes)-1 || suffixes[i]+1 != suffixes[i+1] {
				cidrList, err := mapcidr.IpRangeToCIDR(
					fmt.Sprintf("%s.%d", key, suffixes[startIndex]),
					fmt.Sprintf("%s.%d", key, suffixes[i]),
				)
				if err != nil {
					return nil, fmt.Errorf("parsing cidr: %s, key=%s, suffixes=%v", err, key, suffixes)
				}

				for _, cidr := range cidrList {
					count, _ := mapcidr.AddressCount(cidr)
					cidrMap[cidr] = int(count)
				}

				startIndex = i + 1
			}
		}
	}

	return cidrMap, nil
}
