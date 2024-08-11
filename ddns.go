package main

import (
	"ddns/ip"
	"encoding/json"
	"fmt"
)

func main() {
	ipv6NetInterfaces, err := ip.Ipv6()
	if err != nil {
		fmt.Println("Error getting ipv6 ", err)
	}
	for _, netInterface := range ipv6NetInterfaces {
		n, _ := json.Marshal(netInterface)
		fmt.Println(string(n))
	}
}
