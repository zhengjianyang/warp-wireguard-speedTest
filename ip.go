package main

import (
	"fmt"
	"net"
)

// ParseIPRange 解析CIDR字符串或单个IP字符串，并返回IP地址列表。
func ParseIPRange(ipRange string) ([]net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(ipRange)
	if err == nil {
		// 这是一个CIDR地址
		var ips []net.IP
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			// 为每个IP创建一个新的IP地址切片，以避免修改同一个底层数组。
			newIP := make(net.IP, len(ip))
			copy(newIP, ip)
			ips = append(ips, newIP)
		}
		// 删除网络地址和广播地址
		if len(ips) > 2 {
			return ips[1 : len(ips)-1], nil
		}
		return ips, nil
	}

	// 这不是一个CIDR，尝试解析为单个IP
	ipAddr := net.ParseIP(ipRange)
	if ipAddr == nil {
		return nil, fmt.Errorf("无效的IP地址或CIDR: %s", ipRange)
	}
	return []net.IP{ipAddr}, nil
}

// inc 将IP地址增加1
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
