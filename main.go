package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)





/*
clustersh example.sh --username=xxx --password=xxx

there also should be a `nodes` file in the same directory, in which is node ip in clusters
*/
func main() {

	var username string;
	var password string;
	var ipsfilePath string;

	rootCmd := &cobra.Command{
		Use: "clustersh example.sh",
		Short: "'clustersh example.sh', can run your sh in a cluster, without need to install anything in the cluster",
		Args:cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("username: %s\n", username)
			log.Printf("password: %s\n", password)


			shName := args[0]
			_ = shName
		},
	}


	rootCmd.Flags().StringVarP(&username, "username", "U", "root", "")
	rootCmd.Flags().StringVarP(&password, "password", "P", "root", "")
	rootCmd.Flags().StringVarP(&ipsfilePath, "ips", "I", "ips", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}
