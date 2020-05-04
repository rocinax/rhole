package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/rocinax/rhole/pkg/rhole"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	// *************** Server Default Setting ***************
	viper.SetDefault("ServerPort", 6470)

	// pflag config
	pflag.String("config", "/opt/rocinax/rhole/config", "config: rocinax rhole config directory.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// name of config file (without extension)
	viper.SetConfigName("rhole")
	viper.AddConfigPath(viper.GetString("config"))
	viper.SetConfigType("yaml")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	switch viper.GetString("LogType") {
	case ("Stdout"):
		logrus.SetOutput(os.Stdout)
	case ("File"):
		logFile, err := os.Create(
			path.Join(
				viper.GetString("LogDir"),
				"rhole.log",
			),
		)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(logFile)
	default:
		logrus.SetOutput(os.Stdout)
	}

	switch viper.GetString("LogLevel") {
	case ("Trace"):
		logrus.SetLevel(logrus.TraceLevel)
	case ("Debug"):
		logrus.SetLevel(logrus.DebugLevel)
	case ("Info"):
		logrus.SetLevel(logrus.InfoLevel)
	case ("Warn"):
		logrus.SetLevel(logrus.WarnLevel)
	case ("Error"):
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func main() {

	// port normalize
	var addr string
	if viper.GetInt("ServerPort") <= 1024 || viper.GetInt("ServerPort") > 65535 {
		logrus.WithFields(logrus.Fields{
			"type": "server",
			"app":  "rhole",
		}).Errorf("server port is out of range. :%d", viper.GetInt("ServerPort"))
		panic(fmt.Errorf("ServerPort is out of range. :%d", viper.GetInt("ServerPort")))
	}
	addr = ":" + strconv.Itoa(viper.GetInt("ServerPort"))

	// define and get rigis configuration
	var rhl rhole.Rhole
	err := viper.UnmarshalKey("Rhole", &rhl)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	logrus.WithFields(logrus.Fields{
		"type": "server",
		"app":  "rhole",
	}).Infof("rocinax rhole is started. port: %s", addr)

	// define handle func
	http.HandleFunc("/", rhl.ServeHTTP)

	// run rigis server
	http.ListenAndServe(addr, nil)
}
