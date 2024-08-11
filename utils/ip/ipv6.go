package ip

import (
	"fmt"
	"net"
)

type NetInterface struct {
	Name    string
	Address []string
}

func Ipv6() ([]NetInterface, error) {
	// ipv6 有效地址是 2000::/3
	_, effectiveIpv6, _ := net.ParseCIDR("2000::/3")

	// 获取所有网络接口
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return nil, err
	}

	var ipv6NetInterfaces []NetInterface

	// 遍历网络接口
	for _, netInterface := range allNetInterfaces {
		// 只处理接口状态处于活跃状态的网络接口
		if netInterface.Flags&net.FlagUp == 0 {
			continue
		}

		var ipv6 []string

		// 获取网络接口的地址列表
		if addrs, err := netInterface.Addrs(); err == nil {
			// 遍历网络接口的地址列表
			for _, addr := range addrs {
				// 将接口值addr转换为*net.IPNet类型的指针
				ipNet := addr.(*net.IPNet)
				// 如果子网掩码长度 bits 为 128（即IPv6地址）,如果子网掩码长度 bits 为 32（即IPv4地址）
				_, bits := ipNet.Mask.Size()
				if bits == 128 && effectiveIpv6.Contains(ipNet.IP) {
					ipv6 = append(ipv6, ipNet.IP.String())
				}
			}

			if len(ipv6) > 0 {
				ipv6NetInterfaces = append(ipv6NetInterfaces, NetInterface{
					Name:    netInterface.Name,
					Address: ipv6,
				})
			}
		}
	}
	return ipv6NetInterfaces, nil
}
