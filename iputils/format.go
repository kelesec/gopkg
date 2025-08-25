package iputils

import (
	"net"
	"sort"
	"strings"
)

type Format struct {
	values []string
}

func NewFormat(values []string) *Format {
	for i := 0; i < len(values); i++ {
		values[i] = strings.TrimSpace(values[i])
	}
	return &Format{values: values}
}

// Values 获取 values
func (f *Format) Values() []string {
	return f.values
}

// Deduplicate 对 values 进行去重
func (f *Format) Deduplicate() *Format {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(f.values))

	for _, value := range f.values {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}

	f.values = result
	return f
}

// Sort 对 values 进行排序，按照字符串进行排序
func (f *Format) Sort() *Format {
	if len(f.values) != 0 {
		sort.Strings(f.values)
	}
	return f
}

// Verify 对 values 的每一个元素进行格式验证，比如默认验证元素是否是IP地址格式
// @param callback: 数据格式验证回调，如果提供则使用 callback 进行元素格式验证
func (f *Format) Verify(callback func([]string) []string) *Format {
	if callback != nil {
		f.values = callback(f.values)
		return f
	}

	// 默认进行IP地址格式
	newValues := make([]string, 0, len(f.values))
	for _, value := range f.values {
		if net.ParseIP(value) != nil {
			newValues = append(newValues, value)
		}
	}

	f.values = newValues
	return f
}

// FilterValues 使用 filterFunc 对 values 元素进行过滤
func (f *Format) FilterValues(filterFunc func(value string) bool) []string {
	if filterFunc == nil {
		return f.values
	}

	result := make([]string, 0, len(f.values))
	for _, value := range f.values {
		if filterFunc(value) {
			result = append(result, value)
		}
	}

	return result
}

// FilterPrivateIPs 获取私有IP
func (f *Format) FilterPrivateIPs() []string {
	return f.FilterValues(func(value string) bool {
		if ip := net.ParseIP(value); ip != nil && ip.IsPrivate() {
			return true
		}
		return false
	})
}

// FilterPublicIPs 获取公网IP
func (f *Format) FilterPublicIPs() []string {
	return f.FilterValues(func(value string) bool {
		if ip := net.ParseIP(value); ip != nil && !ip.IsPrivate() {
			return true
		}
		return false
	})
}
