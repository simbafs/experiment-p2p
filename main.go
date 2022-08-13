package main

import (
	"log"
	"os"

	"github.com/simbafs/experiment-p2p/client"
	"github.com/simbafs/experiment-p2p/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gop2p",
	Short: "gop2p is a experiment of p2p implemented in golang",
	Long:  "gop2p is a experiment of p2p implemented in golang",
}

func init() {

	rootCmd.AddCommand(client.Cmd)
	rootCmd.AddCommand(server.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
