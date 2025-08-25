# iputils

## IP 地址生成

- 通过 IP 地址范围生成 IP 地址
```go
package main

import (
	"github.com/kelesec/gopkg/iputils"
	"log"
)

func main() {
	ips, err := iputils.IpsFromRange("192.168.0.10-192.168.0.49")
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		log.Println(ip)
	}

	ips, err = iputils.IpsFromRange("192.168.0.100-126")
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		log.Println(ip)
	}

}
```

- 通过 CIDR 生成 IP 地址
```go
package main

import (
	"github.com/kelesec/gopkg/iputils"
	"log"
)

func main() {
	ips, err := iputils.IPsFromCIDR("192.168.0.1/26")
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		log.Println(ip)
	}
}
```

## CIDR 整理

- 从杂乱的字符串数组中提取 IP_V4 地址，并生成 CIDR
```go
package main

import (
	"fmt"
	"github.com/kelesec/gopkg/iputils"
	"log"
)

func main() {
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

	cidrList, err := iputils.CIDRFromIps(ips)
	if err != nil {
		log.Fatal(err)
	}

	for cidr, count := range cidrList {
		log.Printf("%s: %d\n", cidr, count)
	}
}

```

## IP 地址整理

- 将杂乱的字符串进行去重、进行格式校验、排序、并提取公/私网 IP 地址
```go
package main

import (
	"github.com/kelesec/gopkg/iputils"
	"log"
	"net"
)

func main() {
	values := iputils.NewFormat([]string{
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
	}).Deduplicate(). // 数据去重
		Verify(nil). // 进行IP地址格式验证
		Sort().      // 进行排序
		Values()     // 获取处理后的数据

	for _, v := range values {
		log.Println(v)
	}

	// 获取私有IP地址
	log.Println("Private ips:")
	privateIps := iputils.NewFormat(values).FilterPrivateIPs()
	for _, v := range privateIps {
		log.Println(v)
	}

	// 通过 FilterValues 过滤函数获取公网IP地址
	log.Println("Public ips:")
	publicIps := iputils.NewFormat(values).FilterValues(func(value string) bool {
		if ip := net.ParseIP(value); ip != nil && !ip.IsPrivate() {
			return true
		}
		return false
	})
	for _, v := range publicIps {
		log.Println(v)
	}
}

```

