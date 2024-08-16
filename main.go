package main

import (
	"ddns/utils/dns"
	"ddns/utils/effective"
	"ddns/utils/ip"
	"flag"
	"fmt"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"strings"
	"time"
)

var secretId = flag.String("secretId", "", "SecretId")
var secretKey = flag.String("secretKey", "", "SecretKey")
var domain = flag.String("domain", "", "Domain")
var subdomain = flag.String("subdomain", "", "Subdomain")
var interval = flag.Int64("interval", 30, "Interval")
var subdomains []string

func main() {
	// 解析命令行参数
	flag.Parse()
	if !effective.EffectiveString(*secretId, *secretKey, *domain, *subdomain) {
		fmt.Println("secretId and secretKey and domain and subdomain are required!")
		return
	}
	subdomains = strings.Split(*subdomain, ",")

	// 初始化 dnspod
	dns.InitClient(*secretId, *secretKey)
	// 启动先触发一次
	dnns()

	// 定时查看 ip 是否发生变化
	ticker := time.NewTicker(time.Duration(*interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dnns()
		}
	}
}

func EffectiveRecordAddress(ipv6NetInterfaces []ip.NetInterface, address string) bool {
	f := false
	// 遍历所有 ipv6 地址，查看是否有与 dns 解析一致的地址
	for _, netInterface := range ipv6NetInterfaces {
		for _, addr := range netInterface.Address {
			if addr == address {
				f = true
			}
		}
	}
	return f
}

func dnns() {
	fmt.Println("Execution time: ", time.Now())
	// 查询 dns 记录
	var records []*dnspod.RecordListItem
	for _, sd := range subdomains {
		record := dns.SelectRecord(*domain, sd)
		if record == nil {
			fmt.Printf("No %s's records were found from dnspod!\n", sd)
			panic("Some subdomain has no record!")
		}
		records = append(records, record)
	}
	// 获取当前 ipv6 地址
	ipv6NetInterfaces, err := ip.Ipv6()
	if err != nil {
		fmt.Println("Obtain ipv6 failed!", err)
		return
	}
	for _, record := range records {
		// 如果所有的 ipv6 地址与 dns 解析的地址都不一致，将设备中国第一个 ipv6 地址更新到 dns 解析中，并更新记录值
		if !EffectiveRecordAddress(ipv6NetInterfaces, *record.Value) {
			fmt.Printf("The ipv6 address of record '%s' needs to be updated to '%s'\n", *record.Name, ipv6NetInterfaces[0].Address[0])
			dns.UpdateRecord(*domain, *record.Name, ipv6NetInterfaces[0].Address[0], *record.RecordId)
		} else {
			fmt.Printf("The ipv6 address of record '%s' does not need to be updated\n", *record.Name)
		}
	}
}
