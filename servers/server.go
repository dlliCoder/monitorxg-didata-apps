package servers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"monitorXG/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// var svcconf utils.SvConfig

// 使用配置 启动服务端
func ServerWeb(conf utils.SvConfig) (serverIp string, serverPort string) {

	IP := conf.Server.Host
	Port := conf.Server.Port

	return IP, Port
}

// 定义一个结构体类型用于匹配要接收的JSON数据结构
type Message struct {
	Data string `json:"data"`
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

// 判断函数是否符合标准 暂未使用
func isJsonValidAgainstStruct(msg Message) (bool, error) {
	data := msg.Data
	var dataSt DataStruct
	err := json.Unmarshal([]byte(data), &dataSt)

	if err != nil {
		// 如果解码失败，说明JSON格式不符合结构体
		return false, err
	}

	// 解码成功，可以认为JSON字符串符合结构体（至少在字段名和基本数据类型上匹配）
	return true, nil
}

// 解析msg中的data
func dataAnalyseFunc(msg Message) DataStruct {
	data := msg.Data
	var dataSt DataStruct

	err := json.Unmarshal([]byte(data), &dataSt)
	if err != nil {
		log.Println("解析msg消息JSON时出错: ", err)
		//fmt.Println("解析msg消息JSON时出错 fmt: ", err)
		//return dataSt
	}
	return dataSt
}

// 后续要实现告警发送时，需要实现一个接口使用多样性实现方法重用，目前只有一个邮件发送（调用的时内部提供的jar包）
func SendMail(data DataStruct) {
	// fmt.Println("发送邮件部分，待补充")
	rowler := []string{data.IP, data.Hostname, data.Appname, data.Alive, data.Date, data.DateTime}
	record := strings.Join(rowler, ",")
	log.Println("[I] 邮件发送记录： ", record)

	// 邮件发送待补充
	// 邮件发送补充完成 2024年2月6日 15点42分
	var sendemail utils.JsonContentForEmailSending
	var datasend utils.DataStruct

	// 值传递
	datasend.IP = data.IP
	datasend.Hostname = data.Hostname
	datasend.Appname = data.Appname
	datasend.Alive = data.Alive
	datasend.Date = data.Date
	datasend.DateTime = data.DateTime

	// 赋值发送邮件
	sendemail.WaingMsgSend(datasend)
}

// 告警恢复邮件
func RecoverSendMail(data DataStruct) {
	// fmt.Println("发送邮件部分，待补充")
	rowler := []string{data.IP, data.Hostname, data.Appname, data.Alive, data.Date, data.DateTime}
	record := strings.Join(rowler, ",")
	log.Println("[I] 告警恢复邮件，邮件发送记录： ", record)

	// 邮件发送待补充
	// 邮件发送补充完成 2024年2月6日 15点42分
	var recoverSendemail utils.JsonContentForEmailSending
	var datasend utils.DataStruct

	// 值传递
	datasend.IP = data.IP
	datasend.Hostname = data.Hostname
	datasend.Appname = data.Appname
	datasend.Alive = data.Alive
	datasend.Date = data.Date
	datasend.DateTime = data.DateTime

	// 赋值发送邮件
	recoverSendemail.RecoverMsgSend(datasend)
}

func dataWriter(data DataStruct) {
	// 判断数据有效性
	// 日期为空的字符串不执行
	if data.Date == "" {
		log.Println("没有获取到日期值，不执行写入文件操作")
		// return
	}
	log.Println("将数据转换成csv的格式记录到已日期为文件名的csv文件中")
	rowler := []string{data.IP, data.Hostname, data.Appname, data.Alive, data.Date, data.DateTime}
	record := strings.Join(rowler, ",")
	// 创建一个新的CSV阅读器，从给定的字符串读取数据
	reader := csv.NewReader(strings.NewReader(record))

	filename := data.Date + "-monitor.csv"
	log.Println("filename is :", filename)
	filePath := "../data/" + filename

	// 使用os.Stat()检查文件是否存在
	_, err := os.Stat(filePath)

	if err != nil {
		// 判断错误是否为文件不存在（os.IsNotExist(err)）
		if os.IsNotExist(err) {
			log.Printf("File does not exist : %s\n", filePath)

			// 文件不存在，创建一个新的CSV writer用于写入文件
			file, err := os.Create(filePath)
			if err != nil {
				log.Println("文件不存在，创建文件时报错： ", err)
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush() // 刷新缓存确保所有数据已写入文件

			// 逐行读取并写入文件
			// record, err := reader.Read()
			records, err := reader.ReadAll()

			record := records[0]
			log.Println("record 0 is :", record)
			if err != nil {
				log.Println("文件不存在，创建文件后读取要写入的数据记录报错： ", err)
			}

			err = writer.Write(record)
			if err != nil {
				log.Println("文件不存在，创建文件后并读取要写入的数据记录后写入文件时报错： ", err)
			}
			// reader.
			writer.Flush()
			file.Close()

			log.Println("IsNotExist CSV data has been written to output.csv")

		} else {
			log.Println("Error occurred:", err)
		}
	} else {
		log.Printf("File exists : %s \n", filePath)

		// 打开文件，如果文件已打开并且是可读的，这通常是允许的
		// file, err := os.Open(filePath)
		flag := os.O_APPEND | os.O_WRONLY | os.O_CREATE
		// 设置文件权限，这里是-rwxr-xr-x
		mode := fs.FileMode(0755)

		// 打开文件，即使它已经由其他进程打开用于读取
		file, err := os.OpenFile(filePath, flag, mode)
		// file, err := os.OpenFile(filePath, flag)
		if err != nil {
			log.Println("Error opening file:", err)
			// return
		}

		// defer file.Close()

		info, _ := os.Stat(filePath)
		modes := info.Mode()
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println("Permission denied")
			} else if os.IsExist(err) && (os.FileMode(modes).IsRegular()) { // 检查是否为常规文件且存在
				fmt.Println("File exists but cannot be opened for some reason") // 可能由于其他原因无法打开，如已被其他进程以独占方式打开等
			} else if os.IsNotExist(err) {
				fmt.Println("File does not exist")
			} else {
				fmt.Println("Unknown error:", err)
			}
			return
		}

		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush() // 刷新缓存确保所有数据已写入文件

		// 写入文件
		records, err := reader.ReadAll()
		record := records[0]

		fmt.Println("record is :", record)

		// records, erro := reader.ReadAll()

		if err != nil {
			log.Println("文件存在，读取数据记录时报错： ", err)
		}

		err = writer.Write(record)
		if err != nil {
			log.Println("文件存在，写入数据记录时报错： ", err)
		}

		writer.Flush()
		file.Close()

		log.Println("Exist CSV data has been written to output.csv")
	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 解析请求体中的JSON数据
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received message: %v\n", msg.Data)

	// 可以在此处添加业务逻辑处理...
	// 此处需要添加业务逻辑，处理接收到的数据： data,判断当data： msg.Data，数据打印到本地格式采用csv；
	// 判定data中的一个字段来确认是否需要启动一个邮件服务来发送邮件（待补充）

	data := dataAnalyseFunc(msg)
	fmt.Println("handleRequest : ", data)
	// msg.Data.data
	isalive, _ := strconv.ParseBool(data.Alive)
	if isalive != true {
		// 发送告警邮件
		SendMail(data)
		// 发送邮件后记录数据信息
		dataWriter(data)
	} else {
		// 告警恢复邮件
		RecoverSendMail(data)
		// 记录数据信息
		dataWriter(data)
	}

	// 并向客户端返回响应（例如：成功或错误信息）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func WebStarter(ip string, port string) {
	http.HandleFunc("/", handleRequest)

	log.Printf("[I] Server is listening on port %s...\n\n", port)

	// 拼接监听地址字符串
	var sb strings.Builder
	// sb.WriteString("http://")
	sb.WriteString(ip)
	sb.WriteString(":")
	sb.WriteString(port)
	adress := sb.String()

	log.Printf("[I] sb.string address is : http://%s\n", adress)

	err := http.ListenAndServe(adress, nil)
	// err := http.ListenAndServe("", nil)
	if err != nil {
		log.Println("[E] 启动server段报错： ", err)
		panic(err)
	}
}
