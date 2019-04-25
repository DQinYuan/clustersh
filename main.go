package main

import (
	"fmt"
	"github.com/DQinYuan/clustersh/sshtool"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"log"
     _ "net/http/pprof"
	"net/http"
	"os"
	"runtime"
	"sync"
)


/*
clustersh example.sh --username=xxx --password=xxx

there also should be a `nodes` file in the same directory, in which is node ip in clusters
*/
func main() {
	var username string
	var password string
	var ipsfilePath string
	var timeout string
	var verbose bool
	var debug bool
	var concurrent int
	var command string
	var params string

	rootCmd := &cobra.Command{
		Use: "it can be used in two ways: clustersh shname or clustersh --cmd 'a shell command'",
		Short: "'clustersh shname', can run your sh in a cluster, without need to install anything in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose{
				log.Printf("username: %s\n", username)
				log.Printf("password: %s\n", password)
				log.Printf("ips: %s\n", ipsfilePath)
				log.Printf("timeout: %s\n", timeout)
				log.Printf("verbose: %v\n", verbose)
				log.Printf("concurrent num: %d\n", concurrent)
				log.Printf("cmd: %s", command)
			}

			if debug{
				// start pprof for debug
				go func() {
					log.Println("starting pprof server...")
					err := http.ListenAndServe("localhost:10000", nil)
					if err != nil{
						log.Printf("pprof start error, %v", err)
						return
					}
					log.Println("pprof start success!!!!!")
				}()
			}

			// auth params
			if len(args) < 1 && command == ""{
				fmt.Println("you must give a shname or a cmd")
				fmt.Println(cmd.UsageString())
				os.Exit(1)
			}
			if len(args) >= 1 && command != ""{
				fmt.Println("can not give shname and cmd at the same time")
				fmt.Println(cmd.UsageString())
				os.Exit(1)
			}

			go readNodes(ipsfilePath)

			wg := new(sync.WaitGroup)
			wg.Add(concurrent)

			isSh := command == ""

			var shName string
			var remoteDir string
			var handler func(*sshtool.Sshtool) error
			if isSh{
				shName = args[0]
				uuid := uuid.NewV4().String()
				remoteDir = fmt.Sprintf("~/%s", uuid)
				handler = shHandler(remoteDir, shName, params, verbose)
			} else {
				handler = cmdHandler(command, verbose)
			}

			for i := 0; i < concurrent; i++{
				go clusterExec(handler, username, password, timeout, wg)
			}

			wg.Wait()

			log.Printf("All ok, success num: %d", counter)
		},
	}


	rootCmd.Flags().StringVarP(&username, "username", "U", "root", "")
	rootCmd.Flags().StringVarP(&password, "password", "P", "root", "")
	rootCmd.Flags().StringVarP(&ipsfilePath, "ips", "I", "ips", "")
	rootCmd.Flags().StringVarP(&timeout, "timeout", "T", "10s",  "ssh connect timeout for example:'--timeout 10s'")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "V", false, "print all possible info, default false, can be opened with '--verbose'")
	rootCmd.Flags().BoolVarP(&debug, "debug", "D", false, "start pprof for debug")
	rootCmd.Flags().IntVarP(&concurrent, "concurrent", "C", runtime.NumCPU(), "concurrent num in ssh connection, default value is cpu core num")
	rootCmd.Flags().StringVar(&command, "cmd", "", "a cmd you want to execute in cluster")
	rootCmd.Flags().StringVar(&params, "param", "", "pass params to your script, please wrap them with ''")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}
