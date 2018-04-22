package cmd

import (
	"doko/util"
	"doko/server"
	"fmt"
	"github.com/qiniu/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	web      func()
	role     string
	domain   string
	port     string
	StopChan chan interface{}
)
var rootCmd = &cobra.Command{
	Use:   "doko",
	Short: "doko is a reverse proxy tool",
	Long:  ``,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("out")
	},
}

var runServerCmd = &cobra.Command{
	Use:   "runServer",
	Short: "start the doko server",
	Long:  ``,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		StopChan = util.NewChan(0)
		server.Main(StopChan)
	},
}

//var stopServerCmd = &cobra.Command{
//	Use:   "stopServer",
//	Short: "stop the doko server",
//	Long:  ``,
//	//Args:  cobra.MinimumNArgs(1),
//	Run: func(cmd *cobra.Command, args []string) {
//		StopChan<-1
//	},
//}

var webDCmd = &cobra.Command{
	Use:   "webd",
	Short: "start the doko server web.",
	Long:  ``,
	//Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(port)
		web()
		log.Info("ok")
	},
}



func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "client or server (default is server")

	rootCmd.AddCommand(runServerCmd)
	rootCmd.AddCommand(webDCmd)
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "0.0.0.0", "the public network address of the server")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "4443", "the tunnel port of the server")

}

func initConfig() {
	redis := util.RedisClient()
	err := redis.Set("port", port, 5000*time.Hour).Err()
	if err != nil {
		log.Printf("[initConfig]fail %v", err)
	}
	err = redis.Set("domain", domain, 5000*time.Hour).Err()
	if err != nil {
		log.Printf("[initConfig]fail %v", err)
	}
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
