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
		// 根据网络掩码长度决定是否过滤网络地址和广播地址
		ones, bits := ipNet.Mask.Size()
		if ones < bits-1 {
			// 对于 /30 及以下的网络，过滤掉网络地址和广播地址
			if len(ips) > 2 {
				return ips[1 : len(ips)-1], nil
			}
		}
		// 对于 /31 和 /32 网络，返回所有地址
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
