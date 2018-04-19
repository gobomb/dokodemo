package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"doko/server"
	"log"
)

var web func()

var rootCmd = &cobra.Command{
	Use:   "doko-server",
	Short: "doko is a reverse proxy tool",
	Long:  ``,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "start the doko server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.Main()
	},
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "start the doko server web.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(args)
		web()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(webCmd)
}

func Execute(gin func()) {
	web = gin
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
