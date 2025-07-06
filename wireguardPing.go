package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	pingTimeout                 = time.Millisecond * 1000
	wireguardHandshakeRespBytes = 92
	readBufferSize              = 128 // 必须大于等于 wireguardHandshakeRespBytes
)

// 预定义的默认WireGuard握手包，适用于Cloudflare WARP。
// 这是一个类型1的握手发起消息。
var defaultHandshakePacket, _ = hex.DecodeString("013cbdafb4135cac96a29484d7a0175ab152dd3e59be35049beadf758b8d48af14ca65f25a168934746fe8bc8867b1c17113d71c0fac5c141ef9f35783ffa5357c9871f4a006662b83ad71245a862495376a5fe3b4f2e1f06974d748416670e5f9b086297f652e6dfbf742fbfc63c3d8aeb175a3e9b7582fbc67c77577e4c0b32b05f92900000000000000000000000000000000")

// bufferPool 是一个sync.Pool，用于重用读取缓冲区，以减少内存分配和GC压力。
var bufferPool = sync.Pool{
	New: func() interface{} {
		// 创建一个新的缓冲区，当池中没有可用的缓冲区时调用
		return make([]byte, readBufferSize)
	},
}

// WireGuardPing 向WireGuard端点执行单个握手以测量延迟。
// 它发送一个预定义的握手数据包并等待响应。
// 成功时返回true和往返时间(RTT)，失败时返回false和零持续时间。
func WireGuardPing(ip net.IP, port int) (bool, time.Duration) {
	// 正确格式化地址，支持IPv4和IPv6。
	// 对于IPv6，主机部分必须用方括号括起来。
	var fullAddress string
	if ip.To4() == nil {
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), port)
	} else {
		fullAddress = fmt.Sprintf("%s:%d", ip.String(), port)
	}

	// 建立一个带超时的UDP连接。
	conn, err := net.DialTimeout("udp", fullAddress, pingTimeout)
	if err != nil {
		return false, 0
	}
	defer func() { _ = conn.Close() }()

	// 记录开始时间并发送握手包。
	startTime := time.Now()
	_, err = conn.Write(defaultHandshakePacket)
	if err != nil {
		return false, 0
	}

	// 为读取操作设置截止日期，以避免无限期阻塞。
	err = conn.SetDeadline(time.Now().Add(pingTimeout))
	if err != nil {
		return false, 0
	}

	// 从池中获取一个缓冲区，并确保在函数返回时将其放回池中。
	revBuff := bufferPool.Get().([]byte)
	defer bufferPool.Put(revBuff)

	// 等待并将响应读入缓冲区。
	n, err := conn.Read(revBuff)
	if err != nil {
		return false, 0
	}

	// 一个有效的WireGuard握手响应正好是92字节。
	if n != wireguardHandshakeRespBytes {
		return false, 0
	}

	// 计算并返回持续时间。
	duration := time.Since(startTime)
	return true, duration
}
