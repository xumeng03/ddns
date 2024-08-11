package main

import (
	"ddns/utils/dns"
	"ddns/utils/effective"
	"ddns/utils/ip"
	"flag"
	"fmt"
)

func main() {
	var secretId = flag.String("secretId", "", "SecretId")
	var secretKey = flag.String("secretKey", "", "SecretKey")
	var domain = flag.String("domain", "", "Domain")
	var subdomain = flag.String("subdomain", "", "Subdomain")

	flag.Parse()

	if !effective.EffectiveString(*secretId, *secretKey, *domain, *subdomain) {
		fmt.Println("secretId and secretKey and domain and subdomain are required!")
		return
	}
	dns.InitClient(*secretId, *secretKey)
	record := dns.SelectRecord(*domain, *subdomain)
	if record == nil {
		fmt.Println("No records were found from dnspod!")
		return
	}

	ipv6NetInterfaces, err := ip.Ipv6()
	if err != nil {
		fmt.Println("Error getting ipv6 ", err)
		return
	}
	f := false
	for _, netInterface := range ipv6NetInterfaces {
		for _, addr := range netInterface.Address {
			if addr == *record.Value {
				f = true
			}
		}
	}
	if !f {
		fmt.Println("IPV6 has changed!")
		dns.UpdateRecord(*domain, *subdomain, ipv6NetInterfaces[0].Address[0], *record.RecordId)
	} else {
		fmt.Println("IPV6 has not changed!")
	}
}
