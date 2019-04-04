package main

import (
	"fmt"
	"github.com/DQinYuan/clustersh/sshtool"
	"sync"
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
	t.SkipNow()

	correctData := make([]string, 0)
	correctData = append(correctData, "10.10.108.85")
	correctData = append(correctData, generateIp(`10.10.108.%d`, 91, 93)...)
	correctData = append(correctData, "10.10.108.23")


	go readNodes("ips")
	for _, data := range correctData  {
		ip, _ := <- ch

		if ip != data{
			t.Errorf("Expect:%s, Real: %s ", data, ip)
		}
	}
}


func TestTranAllFiles(t *testing.T) {
	t.SkipNow()
	sshTool, _ := sshtool.NewSshtool("10.10.108.85", "root", "vt1111", "10s")
	tranAllFiles(sshTool, "~/YYSNCN", false)
}

func TestChooseFile(t *testing.T) {
	fileContent = map[string]*CachedFile{
		"example.sh": new(CachedFile),
	}

	choosed := chooseFile("example", "centos")
	if choosed != "example.sh"{
		t.Errorf("Expected: example.sh, Real: %s", choosed)
	}

	fileContent["example_centos.sh"] = new(CachedFile)
	choosed = chooseFile("example", "centos")
	if choosed != "example_centos.sh"{
		t.Errorf("Expected: example_centos.sh, Real: %s", choosed)
	}
}

func TestExecSh(t *testing.T)  {
	t.SkipNow()
	ch <- "10.10.108.40"
	close(ch)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	execSh("~/NNNNNYYY", "test", "root", "vt1111", "5s", true, wg)
	fmt.Println(counter)
}

func TestCount(t *testing.T){

	wg := sync.WaitGroup{}
	wg.Add(2)

	countFunc := func() {
		defer wg.Done()
		for i := 0; i < 30; i++{
			count()
		}
	}

	go countFunc()
	go countFunc()

	wg.Wait()
	if counter != 60{
		t.Errorf("Counter Error, Expected: 60, Real: %d", counter)
	}
}
