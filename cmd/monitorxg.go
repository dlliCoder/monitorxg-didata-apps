package main

import (
	"fmt"
	"log"
	"monitorXG/applications"
	"monitorXG/servers"
	"monitorXG/utils"
)

// 主函数，判断启动的角色来启动不同的服务部分
func main() {
	fmt.Println("hello monitor")
	// 分析参数来判断是启动哪个服务角色
	// 先判断程序执行后面跟的参数：如果没有其他参数跟随
	// cschosen, filename := AnalyseArgs()

	// 分析参数
	cschosen, filename := utils.AnalyseArgs()

	// 配置分析
	switch cschosen {
	case "server":
		// fmt.Println(cschosen)
		config := utils.SvConfig{}
		conf, err := config.ConfigReader(filename)
		if err != nil {
			fmt.Println("server's configer error is :", err)
		} else {
			// 无报错的情况下：角色为server
			fmt.Printf("server's conf struct: %+v\n\n%v\n\n", conf, conf)
			// 分析配置值
			ip, port := servers.ServerWeb(conf)
			servers.WebStarter(ip, port)
		}
	case "appclient":
		config := utils.AppConfig{}
		conf, err := config.ConfigReader(filename)
		if err != nil {
			// err = err
			log.Println("appclietn configer get error: ", err)
		} else {
			// err = nil
			// 无报错的情况下：角色为appclient
			log.Printf("app's conf struct: %+v\n%v", conf, conf)
			// fmt.Println()
			// 启动客户端
			applications.AppclientStarter(conf)
		}
	// 代码可在此处继续扩展
	default:
		fmt.Println("cschosen, error")
	}
}
