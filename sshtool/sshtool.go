package sshtool

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// shell to get system type
const osTypeQuery  = `
Get_Dist_Name()
{
    if grep -Eqii "CentOS" /etc/issue || grep -Eq "CentOS" /etc/*-release; then
        echo 'centos'
    elif grep -Eqi "Red Hat Enterprise Linux Server" /etc/issue || grep -Eq "Red Hat Enterprise Linux Server" /etc/*-release; then
        echo 'rhel'
    elif grep -Eqi "Aliyun" /etc/issue || grep -Eq "Aliyun" /etc/*-release; then
        echo 'aliyun'
    elif grep -Eqi "Fedora" /etc/issue || grep -Eq "Fedora" /etc/*-release; then
        echo 'fedora'
    elif grep -Eqi "Debian" /etc/issue || grep -Eq "Debian" /etc/*-release; then
        echo 'debian'
    elif grep -Eqi "Ubuntu" /etc/issue || grep -Eq "Ubuntu" /etc/*-release; then
        echo 'ubuntu'
    elif grep -Eqi "Raspbian" /etc/issue || grep -Eq "Raspbian" /etc/*-release; then
        echo 'raspbian'
    else
        echo 'unknow'
    fi
}

Get_Dist_Name
`


type Sshtool struct {
	client *ssh.Client
}


func NewSshtool(ip string, username string, passward string, timeout string) (*Sshtool, error) {
	dura, err := time.ParseDuration(timeout)
	if err != nil {
		log.Printf("Warning: %s is not legal time format, use default value '15s'  \n", timeout)
		dura = 15 * time.Second
	}

	config := &ssh.ClientConfig{
		// username
		User: username,
		Auth: []ssh.AuthMethod{
			//password
			ssh.Password(passward),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: dura,
	}
	client, err := ssh.Dial("tcp", ip + ":22", config)
	if err != nil {
		return nil, err
	}

	return &Sshtool{client: client}, nil
}

/*
query linux os distributed version, inclue ubuntu, centos, aliyun, fedora, debian, raspbian
*/
func (p *Sshtool) OsType(verbose bool) (string, error)  {
	oType, err :=  p.Query(osTypeQuery, verbose)
	if err != nil{
		return "", err
	}
	return strings.TrimSpace(oType), nil
}

func (p *Sshtool) Query(cmd string, verbose bool) (string, error) {
	session, err := p.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	dw := &doubeWriter{verbose: verbose}
	session.Stdout = dw
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return dw.string(), nil
}

func (p *Sshtool) Exec(cmd string, verbose bool) error  {
	session, err := p.client.NewSession()
	if err != nil {
		return err
	}
	if verbose{
		session.Stdout = os.Stdout
	}
	if err := session.Run(cmd); err != nil {
		return err
	}

	return nil
}

func (p *Sshtool) Sh(shFilePath string, verbose bool) error  {
	return p.Exec(fmt.Sprintf("sh %s", shFilePath), verbose)
}

func (p *Sshtool) Mkdir(path string, verbose bool) error  {
	return p.Exec("mkdir -p " + path, verbose)
}

func (p *Sshtool) RmDir(remoteDir string, verbose bool) error  {
	return p.Exec("rm -rf " + remoteDir, verbose)
}


func (p *Sshtool) CopyFile(fileReader io.Reader, remotePath string, permissions string, verbose bool) error  {
	contents_bytes, err := ioutil.ReadAll(fileReader)
	if err != nil{
		return err
	}

	bytes_reader := bytes.NewReader(contents_bytes)

	return p.Copy(bytes_reader, remotePath, permissions, int64(len(contents_bytes)), verbose)
}

func (p *Sshtool) Copy(r io.Reader, remotePath string, permissions string, size int64, verbose bool) error {

	filename := path.Base(remotePath)
	directory := path.Dir(remotePath)
	err := p.Mkdir(directory, verbose)
	if err != nil{
		return err
	}

	session, err := p.client.NewSession()
	if err != nil {
		return err
	}
	if verbose{
		session.Stdout = os.Stdout
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()
		w, err := session.StdinPipe()
		if err != nil {
			errCh <- err
			return
		}
		defer w.Close()

		_, err = fmt.Fprintln(w, "C"+permissions, size, filename)
		if err != nil {
			errCh <- err
			return
		}

		// w <- r
		_, err = io.Copy(w, r)
		if err != nil {
			errCh <- err
			return
		}

		//结尾
		_, err = fmt.Fprint(w, "\x00")
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		//scp -t will listen stdin to a file
		err := session.Run("/usr/bin/scp -qt " + directory)
		if err != nil {
			errCh <- err
			return
		}
	}()

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}


func (p *Sshtool) Close() error {
	return p.client.Close()
}


type doubeWriter struct {
	b bytes.Buffer
	verbose bool
}

func (writer *doubeWriter) Write(p []byte) (n int, err error)  {
	if writer.verbose{
		os.Stdout.Write(p)
	}

	writer.b.Write(p)

	return len(p), nil
}

func (writer *doubeWriter) string() string  {
	return writer.b.String()
}