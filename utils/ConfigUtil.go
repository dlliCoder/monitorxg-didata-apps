package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// type returner struct{}

// 接口
type ConfigRead interface {
	// 读取配置
	ConfigReader(file string) (struct{}, error)
	AnalyseConfig(file string) (struct{}, error)
}

/*
	{
	    "server" : {
	        "server_ip" : "127.0.0.1",
	        "server_port": "9900"
	    },
	    "mysql": {
	        "mysqlhost" : "127.0.0.1",
	        "mysql_port": "5230",
	        "mysql_user": "root",
	        "password": "root@123",
	        "database" : "xgmonitor"
	    }
	}
*/

// server端启动配置，server端实例
type SvConfig struct {
	Server   ServerConfig `json:"server"`
	Database DBConfig     `json:"mysql"`
}

type ServerConfig struct {
	Port string `json:"server_port"`
	Host string `json:"server_ip"`
}

type DBConfig struct {
	Mysqlhost string `json:"mysqlhost"`
	Port      string `json:"mysql_port"`
	User      string `json:"mysql_user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
}

func (*SvConfig) ConfigReader(file string) (SvConfig, error) {
	var errre error
	var config SvConfig

	// 判断文件是否存在，存在继续操作并返回错误为nil，不存在返回错误
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file is not exist")
		} else {
			fmt.Println("Err:", err)
		}
		errre = err
	} else {
		// 文件存在
		fmt.Println("file is exist, file info is", fileInfo)
		data, err := os.ReadFile(file)
		if err != nil {
			errre = err
			fmt.Errorf("failed to read config file: %v", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			errre = err
			fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		return config, nil
	}

	return config, errre
}

func (*SvConfig) AnalyseConfig(svcconf SvConfig) map[string]string {
	// var configerMap :=
	configerMap := make(map[string]string)
	return configerMap
}


type AppConfig struct {
	Server       AppServerConfig          `json:"server"`
	Client       AppClientConfig          `json:"client"`
	Applications AppApplicationNameConfig `json:"applicationname"`
}

type AppServerConfig struct {
	ServerPort string `json:"server_port"`
	ServerHost string `json:"server_ip"`
}

type AppClientConfig struct {
	AppPort string `json:"port"`
	AppHost string `json:"ip"`
}

type AppApplicationNameConfig struct {
	Apps string `json:"apps"`
}

// Appconfig
func (*AppConfig) ConfigReader(file string) (AppConfig, error) {
	var errre error
	var config AppConfig

	// 判断文件是否存在，存在继续操作并返回错误为nil，不存在返回错误
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file is not exist")
		} else {
			fmt.Println("Err:", err)
		}
		errre = err
	} else {
		// 文件存在
		fmt.Println("file is exist, file info is", fileInfo)
		data, err := os.ReadFile(file)
		if err != nil {
			errre = err
			fmt.Errorf("failed to read config file: %v", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			errre = err
			fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		return config, nil
	}

	return config, errre
}

// 读取接收者配置文件
func (*Receivers) ConfigReader(file string) (Receivers, error) {

	var errre error
	var config Receivers

	// 判断文件是否存在，存在继续操作并返回错误为nil，不存在返回错误
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file is not exist")
		} else {
			fmt.Println("Err:", err)
		}
		errre = err
	} else {
		// 文件存在
		fmt.Println("file is exist, file info is", fileInfo)
		data, err := os.ReadFile(file)
		if err != nil {
			errre = err
			fmt.Errorf("failed to read config file: %v", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			errre = err
			fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		return config, nil
	}

	return config, errre

}
