package main

import (
	"bytes"
	"fmt"
	"github.com/DQinYuan/clustersh/sshtool"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var ch = make(chan string, runtime.NumCPU())

var (
	// reg to recognize ip like 10.10.108.73
	ipReg = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	// reg to recognize extended ip   like 10.10.108.33-40
	extendedIpReg = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)-(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
)

func readNodes(ipsPath string) {
	data, err := ioutil.ReadFile(ipsPath)
	if err != nil {
		log.Fatalf("%s file can not open\n", ipsPath)
	}

	ips := strings.Split(string(data), "\n")

	for i, v := range ips {
		trimed := strings.TrimSpace(v)
		switch {
		case ipReg.MatchString(trimed):
			// normal ip
			ch <- trimed
		case extendedIpReg.MatchString(trimed):
			// extended ip
			point := strings.LastIndex(trimed, `.`)
			dash := strings.LastIndex(trimed, `-`)
			ipFormat := trimed[:point + 1] + "%d"
			startNum, err := strconv.Atoi(trimed[point + 1:dash])
			if err != nil{
				log.Printf("Warning: line num %d, content: %s, format error\n", i, v)
				continue
			}

			endNum, err := strconv.Atoi(trimed[dash + 1:])
			if err != nil{
				log.Printf("Warning: line num %d, content %s, format error\n", i, v)
				continue
			}

			if endNum < startNum{
				log.Printf("Warning: line num %d, content %s,  start ip can not bigger than end ip", i, v)
			}

			for ip := startNum; ip <= endNum; ip++{
				ch <- fmt.Sprintf(ipFormat, ip)
			}
		}
	}

	close(ch)
}

type CachedFile struct {
	content []byte
	size int64
}

// file relativepath -> file content bytes
var fileContent = make(map[string]* CachedFile)

// get all files in current dir and subdir in memory to fill fileContent map
func init() {
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir(){
			absPath, _ := filepath.Abs(path)
			file, err := os.Open(absPath)
			if err != nil{
				log.Printf("Warning: file %s open fail\n", path)
				return nil
			}

			content, err := ioutil.ReadAll(file)
			if err != nil{
				log.Printf("Warning: file %s read fail", path)
				return nil
			}

			fileContent[path] = &CachedFile{content:content, size:info.Size()}
		}
		return nil
	})

	if err != nil{
		log.Printf("Warning: dir walk error")
	}
}





func shName(remoteDir string, shName string, username string, password string, timeout string, verbose bool) {

	for ip := range ch{
		//create ssh connection
 		sshTool, err := sshtool.NewSshtool(ip, username, password, timeout)
		if err != nil{
			log.Printf("Warning: ip %s can not connect\n", ip)
			continue
		}

		//judge os type
		ostype, err := sshTool.OsType(verbose)
		if err != nil{
			log.Printf("Warning: ip %s os query error\n", ip)
			continue
		}

		//send files in current directory and subdirectory
		for filePath, cFile := range fileContent{
			sshTool.Copy(bytes.NewReader(cFile.content),
				filepath.Join(remoteDir, filePath), "0655", cFile.size, verbose)
		}

		//exec sh for spec os type
	}
}
