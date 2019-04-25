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
	"sync"
	"sync/atomic"
)

var ch = make(chan string, runtime.NumCPU())

var (
	// reg to recognize ip like 10.10.108.73
	//ipReg = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	// reg to recognize extended ip   like 10.10.108.33-40
	extendedIpReg = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)-(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
)

func readNodes(ipsPath string) {
	data, err := ioutil.ReadFile(ipsPath)
	if err != nil {
		log.Fatalf("%s file can not open\n", ipsPath)
	}

	ips := strings.Split(strings.TrimSpace(string(data)), "\n")

	for i, v := range ips {
		trimed := strings.TrimSpace(v)
		switch {
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
		default:
			// normal ip or hostname
			ch <- trimed
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

// transfer all file to remote dir
func tranAllFiles(sshTool *sshtool.Sshtool, remoteDir string, verbose bool){
	for filePath, cFile := range fileContent{
		err := sshTool.Copy(bytes.NewReader(cFile.content),
			filepath.Join(remoteDir, filePath), "0655", cFile.size, verbose)
		if err != nil{
			log.Printf("Warning: copy file error %v", err)
		}
	}
}

func chooseFile(shName string, osType string) string {
	fileLongName := fmt.Sprintf("%s_%s.sh", shName, osType)
	if _, ok := fileContent[fileLongName]; ok{
		return fileLongName
	}

	return fmt.Sprintf("%s.sh", shName)
}

// responsible for ssh connection open and close, counter and WaitGroup
func clusterExec(handler func(*sshtool.Sshtool) error, username, password, timeout string, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range ch{
		log.Printf("start handling ip %s", ip)
		handleIp(ip, username, password, timeout, handler)
	}
}

// responsible for ssh connection open and close, counter
func handleIp(ip, username, password, timeout string, handler func(*sshtool.Sshtool) error) {
	//create ssh connection
	sshTool, err := sshtool.NewSshtool(ip, username, password, timeout)
	if err != nil{
		log.Printf("Warning: ip %s can not connect, err: %v\n", ip, err)
		return
	}
	defer sshTool.Close()

	err = handler(sshTool)
	if err != nil{
		log.Printf("ip %s fail, Warning: %v", ip, err)
		return
	}
	log.Printf("ip %s ok", ip)
	count()
}


func shHandler(remoteDir string, shName, params string, verbose bool) func(*sshtool.Sshtool) error  {

	return func(st *sshtool.Sshtool) error {
		//judge os type
		ostype, err := st.OsType(verbose)
		if err != nil{
			return fmt.Errorf("Warning: os query error, err: %v\n", err)
		}

		//send files in current directory and subdirectory
		tranAllFiles(st, remoteDir, verbose)
		defer st.RmDir(remoteDir, verbose)

		//exec sh for spec os type
		cmd := fmt.Sprintf("cd %s && sh %s %s", remoteDir, chooseFile(shName, ostype), params)
		err = st.Exec(cmd, verbose)
		if err != nil{
			return fmt.Errorf("Warning: %q exec fail, %v \n", cmd, err)
		}

		log.Printf("%q ok\n", cmd)
		return nil
	}
}

func cmdHandler(cmd string, verbose bool) func(*sshtool.Sshtool) error {
	return func(st *sshtool.Sshtool) error {
		err := st.Exec(cmd, verbose)
		if err != nil{
			return fmt.Errorf("Warnning: %q exec fail %v", cmd, err)
		}

		return nil
	}
}


var counter int32

func count() {
	atomic.AddInt32(&counter,1)
}