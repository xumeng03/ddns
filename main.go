package main

import (
	"ddns/utils/dns"
	"ddns/utils/effective"
	"ddns/utils/ip"
	"flag"
	"fmt"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"time"
)

var secretId = flag.String("secretId", "", "SecretId")
var secretKey = flag.String("secretKey", "", "SecretKey")
var domain = flag.String("domain", "", "Domain")
var subdomain = flag.String("subdomain", "", "Subdomain")
var interval = flag.Int64("interval", 15, "Interval")
var record *dnspod.RecordListItem

func main() {
	// 解析命令行参数
	flag.Parse()
	if !effective.EffectiveString(*secretId, *secretKey, *domain, *subdomain) {
		fmt.Println("secretId and secretKey and domain and subdomain are required!")
		return
	}

	// 初始化 dns
	dns.InitClient(*secretId, *secretKey)
	// 查询 dns 记录
	record = dns.SelectRecord(*domain, *subdomain)
	if record == nil {
		fmt.Println("No records were found from dnspod!")
		return
	}
	// 启动先触发一次
	fmt.Println("Execution time: ", time.Now())
	dnns()

	// 定时查看 ip 是否发生变化
	ticker := time.NewTicker(time.Duration(*interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Execution time: ", time.Now())
			dnns()
		}
	}
}

func dnns() {
	// 获取当前 ipv6 地址
	ipv6NetInterfaces, err := ip.Ipv6()
	if err != nil {
		fmt.Println("Error getting ipv6 ", err)
		return
	}
	// 遍历所有 ipv6 地址，查看是否有与 dns 解析一致的地址
	f := false
	for _, netInterface := range ipv6NetInterfaces {
		for _, addr := range netInterface.Address {
			if addr == *record.Value {
				f = true
			}
		}
	}
	// 如果所有的 ipv6 地址与 dns 解析的地址都不一致，将设备中国第一个 ipv6 地址更新到 dns 解析中，并更新记录值
	if !f {
		fmt.Println("IPV6 has changed!")
		dns.UpdateRecord(*domain, *subdomain, ipv6NetInterfaces[0].Address[0], *record.RecordId)
		record = dns.SelectRecord(*domain, *subdomain)
	} else {
		fmt.Println("IPV6 has not changed!")
	}
}
