package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

// 获取进程Name为Java的进程
type AppNamePidList struct {
	Pid             string
	Appname         string
	Cmdline         string
	ApplicationName string
}

type ErrMsg struct {
	Msg string
}

// 实现error接口，自定义一个error
func (e *ErrMsg) Error() string {
	return fmt.Sprintf("error msg : %s", e.Msg)
}

// 进程检查，检查一个appname
func ProcessCheck(appName string) ([]AppNamePidList, error) {
	// 获取所有进程的信息
	processes, err := process.Processes()

	if err != nil {
		log.Println("Error getting processes:", err)
		// return
	}

	var (
		appNameandPidList []AppNamePidList
		appNameandPid     AppNamePidList
		pid               string
		piderr            error
	)

	for _, proc := range processes {
		// 修改匹配名称20240205
		if name, _ := proc.Name(); name == "java" {

			log.Println("process name is name: ", name)

			// 返回进程命令执行信息
			cmdline, _ := proc.Cmdline()

			// proc.

			log.Println("[I] cmdline is : ", cmdline)

			if strings.Contains(cmdline, appName) {
				pid = strconv.FormatInt(int64(proc.Pid), 10)
				appNameandPid.Appname = name
				appNameandPid.Pid = pid
				appNameandPid.Cmdline = cmdline
				appNameandPid.ApplicationName = appName

				appNameandPidList = append(appNameandPidList, appNameandPid)
			}
		}
	}

	// 判断列表长度,来设置报错返回值是否为空
	length := len(appNameandPidList)
	if length == 0 {
		log.Println("appNameandPidList length is 0")
		// 自定义报错信息
		errorMsg := "ProcessCheck appNameandPidList length is 0"
		e := &ErrMsg{
			Msg: errorMsg,
		}
		piderr = e
	} else {
		log.Println("appNameandPidList length is :", length)

	}

	// 返回值
	return appNameandPidList, piderr
}

// AppSocketCheck()  返回可检查的接口列表,待实现

// type AppnameAndPort struct {
// 	Appname string
// 	AppPort string
// }

func AppSocketVerify(list []AppNamePidList, portlists []string) ([]AppnameAndPort, error) {
	// client port 在实际检查中竟没有使用到，发送消息没有使用到端口相关的配置

	/*
		type AppNamePidList struct {
		Pid     string
		Appname string
		Cmdline string
		}

		需要考虑没有获取进程在的情况
		当进程死亡，此方法获取不到进程的id，只能返回其中存活的进程端口
	*/

	var returnDataList []AppnameAndPort
	var returnData AppnameAndPort
	var errorOfAppSocketVerify error

	// 获取pid，通过pid获取程序监听的端口
	for _, data := range list {
		pid := data.Pid
		intPid, errstrconv := strconv.ParseInt(pid, 10, 64)
		if errstrconv != nil {
			log.Println("string conv pid string to int get error: ", errstrconv)
		}

		portList, err := GetPortsByPID(int(intPid))
		if err != nil {
			log.Println("GetPortByPID runs error: ", err)
		}
		portlistStrings := strings.Join(portList, " ")

		log.Println("[I] {AppSocketVerify} data is : ", data, "portlistStrings is ", portlistStrings, "portlists is :", portlists)
		// 端口列表字符串话来判断是否存在端口
		for _, cport := range portlists {
			if strings.Contains(portlistStrings, cport) {
				log.Println("[I] strings.Contains(portlistStrings, cport) :", cport)
				returnData.Appname = data.ApplicationName
				returnData.AppPort = cport
				returnDataList = append(returnDataList, returnData)
				// returnData = append(returnData, cport)
			} else {
				log.Println("returnDataList = append(returnDataList, returnData)")
				continue

			}
		}
	}

	if len(returnDataList) == 0 {
		msg := "port get error in AppSocketVerify"
		log.Println("[E] returnDataList length is : ", len(returnDataList))
		ermsg := ErrMsg{
			Msg: msg,
		}
		errorOfAppSocketVerify = &ermsg
	}

	log.Println("[I] returnDataList is : ", returnDataList)
	log.Println("[I] errorOfAppSocketVerify is : ", errorOfAppSocketVerify)

	return returnDataList, errorOfAppSocketVerify
}

func checkPortReachable(ip string, port int) (bool, error) {
	// 设置超时时间（例如5秒）
	connTimeout := 5 * time.Second

	// 创建TCP连接尝试（可以替换为"udp"以测试UDP端口）
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, connTimeout)

	if err != nil {
		// 如果无法建立连接，则返回错误信息并设置连通性为false
		return false, err
	}

	// 若能成功建立连接，则关闭连接并返回true
	defer conn.Close()
	return true, nil
}

// 将data解析以下字段
// {"ip":"", "hostname":"", "appname":"", "isaive":"", "date":"", "time":""}
type DataStruct struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Appname  string `json:"appname"`
	Alive    string `json:"isaive"`
	Date     string `json:"date"`
	DateTime string `json:"time"`
}

// 传入数据是可以检查的端口号,返回列表,Appn name 也需要一并传输过来，和IP端口对应
func AppSocketCheck(data []AppnameAndPort) []DataStruct {
	var jsondataList []DataStruct
	var jsondata DataStruct

	// 获取主机名称
	hostname, erros := os.Hostname()
	if erros != nil {
		log.Println("[E] get hostname err :", erros)
	}

	// 获取当前时间
	now := time.Now()

	date := now.Format("2006-01-02")
	time := now.Format("15:04:05")

	// 获取Appname
	// appname := ClientAppNameOfAppclient

	// app客户端的ip
	appclientip := ClientIpOfAppClient
	// app客户端的端口
	// appclientport := utils.ClientPortOfAppClient
	jsondata.IP = appclientip
	jsondata.Hostname = hostname
	jsondata.Date = date
	jsondata.DateTime = time

	// 程序存活的使用以下发送数据
	// 包括进程存活端口存在，端口联通 和 端口不联通
	for _, appandport := range data {
		port := appandport.AppPort
		appname := appandport.Appname
		portin, errs := strconv.ParseInt(port, 10, 64)
		if errs != nil {
			log.Println("trans format string port to  int get error : ", errs)
		}

		bl, err := checkPortReachable(appclientip, int(portin))
		if err != nil {
			log.Println("[E] checkPortReachable get err :", err)
		}

		if bl == true {
			log.Println("[I] 生成报送数据")
			// 生成报送数据，追加到报送数据列表（待补充）

			jsondata.Appname = appname
			jsondata.Alive = "true"

			jsondataList = append(jsondataList, jsondata)
		} else {
			log.Println("[E] check port un reachable port is : ", port)
			log.Println("[E] 检查端口不通过，生成报错数据")
			jsondata.Appname = appname
			jsondata.Alive = "false"
			jsondataList = append(jsondataList, jsondata)
		}
	}

	// 以下部分包括 进程不活
	// 当进程不存货时，将不会获取到端口
	// 当进程存活时，获取不到端口的话会返回false

	return jsondataList
}
