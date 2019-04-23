package main

import (
	"fmt"
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
	var username string;
	var password string;
	var ipsfilePath string;
	var timeout string;
	var verbose bool;
	var debug bool;
	var concurrent int;

	rootCmd := &cobra.Command{
		Use: "clustersh shname",
		Short: "'clustersh shname', can run your sh in a cluster, without need to install anything in the cluster",
		Args:cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if verbose{
				log.Printf("username: %s\n", username)
				log.Printf("password: %s\n", password)
				log.Printf("ips: %s\n", ipsfilePath)
				log.Printf("timeout: %s\n", timeout)
				log.Printf("verbose: %v\n", verbose)
				log.Printf("concurrent num: %d\n", concurrent)
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

			shName := args[0]

			uuid := uuid.NewV4().String()
			remoteDir := fmt.Sprintf("~/%s", uuid)

			go readNodes(ipsfilePath)

			wg := new(sync.WaitGroup)
			wg.Add(runtime.NumCPU())

			for i := 0; i < concurrent; i++{
				go execSh(remoteDir, shName, username, password, timeout, verbose, wg)
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}
