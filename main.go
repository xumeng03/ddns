package main

import (
	"ddns/utils/dns"
	"ddns/utils/util"
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
var source = flag.String("source", "netinterface", "Source")
var ipv4 = flag.Bool("ipv4", true, "Ipv4")
var ipv6 = flag.Bool("ipv6", false, "Ipv6")
var interval = flag.Int64("interval", 30, "Interval")

var subdomains []string

func main() {
	// 解析命令行参数
	flag.Parse()
	if !util.Required(*secretId, *secretKey, *domain, *subdomain) {
		fmt.Println("secretId and secretKey and domain and subdomain are required!")
		return
	}
	if !util.Mutex(*ipv4, *ipv6) {
		fmt.Println("IPv4 and ipv6 cannot be turned on at the same time!")
		return
	}
	subdomains = strings.Split(*subdomain, ",")

	// 初始化 dns 配置
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
	fmt.Println(util.WriteValueAsString(records))
	// 获取当前 ip 地址
	var ipInterfaces []util.IPInterface
	var err error
	var recordType string
	if *ipv4 {
		recordType = "A"
		if *source == "netinterface" {
			ipInterfaces, err = util.NetInterfaceIPv4()
			if err != nil {
				fmt.Printf("Failed to obtain ipv4!\n")
				return
			}
		} else if *source == "api" {
			ipInterfaces, err = util.ApiIPv4()
			if err != nil {
				fmt.Printf("Failed to obtain ipv4!\n")
				return
			}
		}
	} else if *ipv6 {
		recordType = "AAAA"
		if *source == "netinterface" {
			ipInterfaces, err = util.NetInterfaceIPv6()
			if err != nil {
				fmt.Printf("Failed to obtain ipv6!\n")
				return
			}
		} else if *source == "api" {
			ipInterfaces, err = util.ApiIPv6()
			if err != nil {
				fmt.Printf("Failed to obtain ipv6!\n")
				return
			}
		}
	}
	fmt.Println(util.WriteValueAsString(ipInterfaces))
	for i := range records {
		record := records[i]
		f := false
		for _, ipInterface := range ipInterfaces {
			if ipInterface.Address == *record.Value {
				fmt.Printf("The ip address of record '%s' does not need to be updated\n", *record.Name)
				f = true
				break
			}
		}
		if !f {
			fmt.Printf("The ip address of record '%s' need to be updated\n", *record.Name)
			dns.UpdateRecord(*domain, *record.Name, recordType, ipInterfaces[0].Address, *record.RecordId)
		}
	}
}
