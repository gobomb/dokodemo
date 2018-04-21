package cmd

import (
	"doko/client"
	"doko/server"
	"fmt"
	"github.com/qiniu/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	web  func()
	role string
)
var rootCmd = &cobra.Command{
	Use:   "doko",
	Short: "doko is a reverse proxy tool",
	Long:  ``,
}

var runCmd = &cobra.Command{
	Use:   "runS",
	Short: "start the doko server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//if viper.config
		server.Main()
	},
}

var runClientCmd = &cobra.Command{
	Use:   "runC",
	Short: "start the doko client",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		client.Main()

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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "client or server (default is server")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(webCmd)
}

func initConfig() {
	if role != "" {
		// Use config file from the flag.
		viper.SetConfigFile(role)
		log.Info(role)
	} else {
		viper.SetConfigFile("server")
	}
}

func Execute(gin func()) {
	web = gin
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
