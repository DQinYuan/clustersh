package clustersh

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

/*
clustersh example.sh

there also should be a `nodes` file in the same directory, in which is node ip in clusters
*/
func main() {
	rootCmd := &cobra.Command{
		Use: "clustersh example.sh",
		Short: "run your sh in a cluster, without need to install anything in the cluster",
		Args:cobra.MinimumNArgs(1),
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}
