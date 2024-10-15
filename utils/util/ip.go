package util

import (
	"fmt"
	"io"
	"net"
)

var v4Apis = []string{"https://ipv4.ddnspod.com"}

var v6Apis = []string{"https://ipv6.ddnspod.com"}

type IPInterface struct {
	Name    string
	Address string
}

func NetInterfaceIPv4() ([]IPInterface, error) {
	// 获取所有网络接口
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return nil, err
	}

	var ipInterfaces []IPInterface

	// 遍历网络接口
	for _, netInterface := range allNetInterfaces {
		// 只处理接口状态处于活跃状态的网络接口
		if netInterface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取网络接口的地址列表
		if addrs, err := netInterface.Addrs(); err == nil {
			// 遍历网络接口的地址列表
			for _, addr := range addrs {
				// 将接口值addr转换为*net.IPNet类型的指针
				ipNet := addr.(*net.IPNet)
				// 如果子网掩码长度 bits 为 128（即IPv6地址）,如果子网掩码长度 bits 为 32（即IPv4地址）
				_, bits := ipNet.Mask.Size()
				if bits == 32 {
					fmt.Println(netInterface.Name, ": ", ipNet.IP.String())
					ipInterfaces = append(ipInterfaces, IPInterface{
						Name:    netInterface.Name,
						Address: ipNet.IP.String(),
					})
				}
			}
		}
	}
	return ipInterfaces, nil
}

func NetInterfaceIPv6() ([]IPInterface, error) {
	// ipv6 有效地址是 2000::/3
	_, effectiveIpv6, _ := net.ParseCIDR("2000::/3")

	// 获取所有网络接口
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return nil, err
	}

	var ipInterfaces []IPInterface

	// 遍历网络接口
	for _, netInterface := range allNetInterfaces {
		// 只处理接口状态处于活跃状态的网络接口
		if netInterface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取网络接口的地址列表
		if addrs, err := netInterface.Addrs(); err == nil {
			// 遍历网络接口的地址列表
			for _, addr := range addrs {
				// 将接口值addr转换为*net.IPNet类型的指针
				ipNet := addr.(*net.IPNet)
				// 如果子网掩码长度 bits 为 128（即IPv6地址）,如果子网掩码长度 bits 为 32（即IPv4地址）
				_, bits := ipNet.Mask.Size()
				if bits == 128 && effectiveIpv6.Contains(ipNet.IP) {
					fmt.Println(netInterface.Name, ": ", ipNet.IP.String())
					ipInterfaces = append(ipInterfaces, IPInterface{
						Name:    netInterface.Name,
						Address: ipNet.IP.String(),
					})
				}
			}
		}
	}
	return ipInterfaces, nil
}

func ApiIPv4() ([]IPInterface, error) {
	client := CreateNoProxyHTTPClient("tcp4")
	var ipInterfaces []IPInterface
	for i := range v4Apis {
		api := v4Apis[i]
		resp, err := client.Get(api)
		if err != nil {
			fmt.Printf("通过接口获取IPv4失败! 接口地址: %s", api)
			continue
		}
		defer resp.Body.Close()
		lr := io.LimitReader(resp.Body, 1024000)
		ip, err := io.ReadAll(lr)
		if err != nil {
			fmt.Printf("通过接口获取IPv4失败! 接口地址: %s", api)
			continue
		}
		ipInterfaces = append(ipInterfaces, IPInterface{api, string(ip)})
	}
	return ipInterfaces, nil
}

func ApiIPv6() ([]IPInterface, error) {
	client := CreateNoProxyHTTPClient("tcp6")
	var ipInterfaces []IPInterface
	for i := range v6Apis {
		api := v6Apis[i]
		resp, err := client.Get(api)
		if err != nil {
			fmt.Printf("通过接口获取IPv6失败! 接口地址: %s", api)
			continue
		}
		defer resp.Body.Close()
		lr := io.LimitReader(resp.Body, 1024000)
		ip, err := io.ReadAll(lr)
		if err != nil {
			fmt.Printf("通过接口获取IPv6失败! 接口地址: %s", api)
			continue
		}
		ipInterfaces = append(ipInterfaces, IPInterface{api, string(ip)})
	}
	return ipInterfaces, nil
}
