package main

import (
	"fmt"
	"rnctl/pkg/zerotier"
)

func main() {
	info, err := zerotier.GetInfo()
	if err != nil {
		fmt.Printf("failed to get status: %v\n", err)
		return
	}
	fmt.Printf("current status: %s\n", info)

	nwid := "我懒得再调用一次vault.go了，就这样吧，反正所有test都是gemini帮忙写的"
	res, err := zerotier.JoinNetwork(nwid)
	if err != nil {
		fmt.Printf("failed to join network: %v\n", err)
	} else {
		fmt.Printf("took %s to join the network\n", res)
	}
	list, _ := zerotier.ListNetworks()
	fmt.Printf("list of the network (JSON): %s\n", list)
}
