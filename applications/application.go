package applications

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"monitorXG/utils"
	"net/http"
	"strings"
	"time"
)

// var appconf utils.AppConfig


// app server报送数据server端配置
func AppServerWebConfig(appconfig utils.AppConfig) (serverIp string, serverPort string) {
	serversIP := appconfig.Server.ServerHost
	serversPort := appconfig.Server.ServerPort
	return serversIP, serversPort
}

// app client报送数据本地启动服务的IP地址s
func AppClientWebConfig(appconfig utils.AppConfig) (clientIP string, clientPort string) {
	clientip := appconfig.Client.AppHost
	clinetport := appconfig.Client.AppPort
	return clientip, clinetport
}

// app client监听的服务和端口信息
func AppClientApplications(appconfig utils.AppConfig) (appsurl string) {
	apps := appconfig.Applications.Apps
	return apps
}

// 同样的Message结构体定义，确保两端一致
type Message struct {
	Data string `json:"data"`
}

// 传送数据为json格式的信息，实际发送数据核心主体方法
func SendJSONToServer(serverip string, severport string, data utils.DataStruct) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// 如果序列化过程中发生错误，这里会捕获并处理它
		log.Println("数据格式转换失败： ", err)
	}

	jsonString := string(jsonBytes)

	msgToSend := Message{
		Data: jsonString,
	}

	// 转换为json的字节流
	jsonData, err := json.Marshal(msgToSend)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return err
	}

	// 拼接字符串server端信息
	var sbs strings.Builder
	sbs.WriteString("http://")
	sbs.WriteString(serverip)
	sbs.WriteString(":")
	sbs.WriteString(severport)
	adressSv := sbs.String()

	// 创建HTTP POST请求
	url := adressSv
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("http.NewReques get error: ", err)
	}
	// 设置Content-Type头信息以表明发送的是JSON数据
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request,  client.Do(req) get error :", err)
		return err
	}
	defer resp.Body.Close()

	// 如果需要，读取并处理响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	} else {
		log.Println("Response Body:", string(body))
	}

	log.Println("Sent JSON data to the server.")
	return nil
}

func ApplicationCheck(apps string) ([]utils.DataStruct, error) {
	var data []utils.DataStruct
	var errOfApplicationCheck error
	// var PostDataFalseValue utils.DataStruct

	// 按照appuri去判断
	appuris := strings.Split(apps, ";")
	for _, appuri := range appuris {
		appNamePorts := strings.Split(appuri, ":")
		appname := appNamePorts[0]

		// 需要检查的端口，端口列表
		ports := appNamePorts[1]
		portList := strings.Split(ports, ",")

		// 检查进程，返回一个appinfo的列表，因为一个app可能匹配多个的情况
		// 当返回的appinfo是一个空列表的话，说明这个app检查的时候没有获取到进程，应当追加列表进去一个值说明这个进程不存在，用来发送错误消息
		appInfoList, err := utils.ProcessCheck(appname)
		if err != nil {
			log.Println("[E] get error in use utils.ProcessCheck: ", err)
			// return data
			// msg := utils.ErrMsg{
			// 	Msg: "[E] get ProcessCheck's error.",
			// }
			// errOfApplicationCheck = &msg
		}

		// 列表长度为0，没有获取到进程，即认为进程没有存活，传送一个发送邮件的消息值给另一个函数
		// 进程不存活的情况，直接跳过端口检查，生成告警消息发发送到服务端

		// 分类判断，长度为0 发送错误数据，长度不为0 进行端口验证，返回数值发送正确数据
		if len(appInfoList) == 0 {
			log.Println("[I] 没有获取到进程: ", appname)
			// PostDataFalseValue.Alive = "false"
			posetFalseData := utils.ApplicationDiedlineProcessMsg(appname)
			data = append(data, posetFalseData)
		} else {
			// 将appinfo传输给检查端口，在配置文件中的端口与appinfo的pid获取的端口中有相同值即可；然后进程端口通信检查
			// 所有端口都存在且进行端口检查通过之后返回正常的data值进行发送

			// var appUriChecking []string

			// 返回需要进行socket端口探活验证的列表，传参给端口验证，返回数据进行发送使用：
			// 返回的发送值列表保存，然后返回
			// 端口列表 是 实际配置文件中的
			portslistVerified, err := utils.AppSocketVerify(appInfoList, portList)
			if err != nil {
				log.Println(err)
				// log.Println("此次循环不执行后续操作，进入")
				// continue
				msg := utils.ErrMsg{
					Msg: "[E] get AppSocketVerify's error.",
				}
				errOfApplicationCheck = &msg
			}

			//  App端口检查步骤，生成json格式的数据返回并用于数据的发送
			// 传入数为字符串列表，返回的是数据列表，但是在for循环中
			postDataList := utils.AppSocketCheck(portslistVerified)

			for _, postDataa := range postDataList {
				data = append(data, postDataa)
			}

		}

	}

	return data, errOfApplicationCheck
}

// 发送数据到server端
func ApplicationDataSend(data []utils.DataStruct, serverIp string, serverPort string) (string, error) {
	var err error
	// 发送数据的逻辑
	for _, cleientData := range data {
		err = SendJSONToServer(serverIp, serverPort, cleientData)
		if err != nil {
			log.Println(err)
		}
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// 如果序列化过程中发生错误，这里会捕获并处理它
		log.Println("数据格式转换失败： ", err)
	}

	jsonString := string(jsonBytes)
	result := jsonString

	return result, err
}

func AppclientStarter(appconfig utils.AppConfig) {
	// 服务端配置
	serverip, serverport := AppServerWebConfig(appconfig)

	// 客户端配置
	clientip, clientport := AppClientWebConfig(appconfig)
	log.Printf("client ip is : %s, client port is : %s", clientip, clientport)
	utils.ClientIpOfAppClient = clientip
	utils.ClientPortOfAppClient = clientport

	//监控的进程名称和端口
	apps := AppClientApplications(appconfig)
	// 分别进行进程检测和端口检测

	// 使用for循环使程序保持持续运行
	for {
		// 每30s检查一次
		log.Println("[I] application sleep 30 second.")
		time.Sleep(time.Second * 30)

		log.Println("[I] application client start.")
		log.Println("[I] application checking...")
		log.Println("[I] application end checking...")

		// for循环 每30秒运行检查
		data, errOfApplicationCheck := ApplicationCheck(apps)
		if errOfApplicationCheck != nil {
			log.Println(errOfApplicationCheck)
			continue
		}
		// 发送数据,程序的核心功能
		result, err := ApplicationDataSend(data, serverip, serverport)

		if err != nil {
			log.Println("[E]data send error, err is :", err, "continue next check.")
			continue
		}
		log.Printf("data send succsess: %s ", result)
	}

}
