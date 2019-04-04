package sshtool

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSshtool(t *testing.T) {
	t.SkipNow()

	username := "root"

	sshTool, err := NewSshtool("10.10.108.85", username, "vt1111", "5s")
	if err != nil{
		t.Errorf("ssh create fail")
	}

	result, err := sshTool.Query("/usr/bin/whoami", false)
	if err != nil{
		t.Errorf("cmd exec error")
	}

	if strings.TrimSpace(result) != "root"{
		t.Errorf("Expected: %s, Real: %s", username, result)
	}

	osType, err := sshTool.OsType(false)
	if err != nil{
		t.Errorf("os query error")
	}

	if strings.TrimSpace(osType) != "centos"{
		t.Errorf("Expected: CentOS, Real: %s", osType)
	}

	absPath, _ := filepath.Abs("sshtool.go")
	file, err := os.Open(absPath)
	if err != nil{
		log.Fatal(err)
	}

	sshTool.CopyFile(file, "/root/clusershtest/sshtool.go", "0655", true)

	sshTool.Mkdir("~/hahahaha", true)
	sshTool.RmDir("~/hahahaha", true)
}
