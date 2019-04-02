package main

import (
	"fmt"
	"testing"
)

func generateIp(ipFormat string, start int, end int) []string {
	ips := make([]string, 0)

	for i := start; i <= end; i++{
		ips = append(ips, fmt.Sprintf(ipFormat, i))
	}

	return ips
}


func TestReadNodes(t *testing.T) {
	correctData := make([]string, 0)
	correctData = append(correctData, "10.10.108.34")
	correctData = append(correctData, "10.10.108.66")
	correctData = append(correctData, generateIp(`10.10.108.%d`, 99, 104)...)
	correctData = append(correctData, "10.10.108.23")


	go readNodes("ips")
	for _, data := range correctData  {
		ip, _ := <- ch

		if ip != data{
			t.Errorf("Expect:%s, Real: %s ", data, ip)
		}
	}
}
