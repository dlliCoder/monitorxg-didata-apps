package utils

import (
	"log"
	"os"
	"time"
)

func ApplicationDiedlineProcessMsg(appname string) DataStruct {
	// 返回数据结构体

	var data DataStruct
	// 获取主机名称
	host, err := os.Hostname()
	if err != nil {
		log.Println("Get HostName error : ", err)
	}
	// HostName = host
	// return host

	// 获取当前时间
	now := time.Now()

	date := now.Format("2006-01-02")
	time := now.Format("15:04:05")

	// 获取Appname
	// appname := ClientAppNameOfAppclient

	// app客户端的ip
	appclientip := ClientIpOfAppClient

	data.Hostname = host
	data.Alive = "false"
	data.Appname = appname
	data.Date = date
	data.DateTime = time
	data.IP = appclientip

	return data
}
