package utils

import (
	"flag"
	"fmt"
	"os"
)

// 打印用法
func flagUsage() {
	usageText := `monitorxg is in mode server/client.
	
	Usage:
	monitorxg command [arguments]
	The command are::
		sevrer: monitor server a string
		appclient: appmonitor client a strng
	Use monitorxg [command] --help`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

// 分析参数
func AnalyseArgs() (cschosen string, filename string) {
	flag.Usage = flagUsage

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	appclientCmd := flag.NewFlagSet("appclient", flag.ExitOnError)

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	var file string
	var cschose string

	switch os.Args[1] {
	case "server":
		f := serverCmd.String("f", "../conf/server.json", "choose sevrer mode to start, f is server mode config file.")
		serverCmd.Parse(os.Args[2:])
		file = *f
		cschose = "server"
	case "appclient":
		f := appclientCmd.String("f", "../conf/application.json", "choose appclient mode to start, f is appclient mode config file")
		appclientCmd.Parse(os.Args[2:])
		file = *f
		cschose = "appclient"
	default:
		flag.Usage()
	}

	// 默认值设置为server
	// cschose := flag.String("cschose", "e-server", "chosen server or client")
	// file := flag.String("f", "e-server.conf", "chosen server or client")

	// flag.Parse()
	// return *cschose, *file
	// 返回选择的模式和配置文件内容
	return cschose, file
}
