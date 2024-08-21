package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	// "github.com/develar/go-fs-util"
)

/*
实现邮件发送告警，使用接口形式
*/

// // {"ip":"", "hostname":"", "appname":"", "isaive":"", "date":"", "time":""}
// type DataStruct struct {
// 	IP       string `json:"ip"`
// 	Hostname string `json:"hostname"`
// 	Appname  string `json:"appname"`
// 	Alive    string `json:"isaive"`
// 	Date     string `json:"date"`
// 	DateTime string `json:"time"`
// }

type WarnigSend interface {
	WaingMsgSend(msg *struct{})
	RecoverMsgSend(msg *struct{})
}

type Receivers struct {
	Receivers string `json:"receivers"`
}

// 使用此结构体实现发送消息的方法
type JsonContentForEmailSending struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Receivers string `json:"receivers"`
}

// var msgdata DataStruct

// 新增告警检查
/*
告警检查WarningCheck():
	1、每次发送告警时检查目录下对应的文件是否存在，文件命令格式--> msgstruct DataStruct 的IP-APPNAME-DATE.txt
	2、文件存在时，不执行下一步的告警信息处理和发送，并读取文件中的数字，按照数字进行静默；
	3、静默后将文件删除
*/
func WarningCheck(msgstruct DataStruct) bool {
	log.Println("[I] 开始处理静默文件")

	var checkResult bool

	filename := msgstruct.IP + msgstruct.Appname + msgstruct.Date + "txt"
	filepath := "../data/warning/" + filename

	// 检查文件是否存在
	_, errosStat := os.Stat(filepath)

	if errosStat != nil {
		// 判断错误是否为文件不存在（os.IsNotExist(err)）
		if os.IsNotExist(errosStat) {
			// 文件不存在，进行下一步
			log.Println("[I] 静默文件不存在，进行发送告警并出进行创建文件...")
			checkResult = false
		} else {
			// 文件存在，读取数据,进行静默
			log.Println("[I]  静默文件存在，进行告警静默10分钟")
			checkResult = true
		}
	} else {
		// 为空说明文件存在
		log.Println("[E] 检查静默文件是否存在时,错误为： ", errosStat)
		// 默认为false，执行发送告警并创建
		checkResult = true

	}

	return checkResult
}

func WarningETL(msgstruct DataStruct) []byte {
	// 警信息处理 ------------------------------------------------------------------------------------------------------------------------------
	// var err ErrMsg
	log.Println("[I] 开始告警信息处理... ...")
	var wmsg JsonContentForEmailSending
	var receiver Receivers

	// 读取配置文件中的接收者
	receiverJson, errreader := receiver.ConfigReader("../conf/emailreciever.json")
	if errreader != nil {
		log.Println("[E] read email config fail: ", errreader)
	}

	ip := msgstruct.IP
	hostname := msgstruct.Hostname
	appname := msgstruct.Appname
	date := msgstruct.Date
	time := msgstruct.DateTime
	isalive := msgstruct.Alive

	aliveStatus, err := strconv.ParseBool(isalive)
	if err != nil {
		log.Println("[E] 转换凭他存活状态类型为bool类型时报错 : ", err)
	}

	// 拼接消息串
	if aliveStatus == true {
		title := []string{"DI.平台服务告警 恢复--", appname}
		titileString := strings.Join(title, "")
		msg := []string{"IP地址: ", ip, "<br>主机名称：", hostname, "<br>", "DI.平台服务检查进程 恢复：", appname, "<br>告警时间：", date, " ", time}
		msgString := strings.Join(msg, "")
		receivers := receiverJson.Receivers

		wmsg.Title = titileString
		wmsg.Content = msgString
		wmsg.Receivers = receivers
	} else {
		title := []string{"DI.平台服务告警--", appname}
		titileString := strings.Join(title, "")
		msg := []string{"IP地址: ", ip, "<br>主机名称：", hostname, "<br>", "DI.平台服务检查进程存活信息失败：", appname, "<br>告警时间：", date, " ", time}
		msgString := strings.Join(msg, "")
		receivers := receiverJson.Receivers

		wmsg.Title = titileString
		wmsg.Content = msgString
		wmsg.Receivers = receivers
	}
	// title := []string{"DI.平台服务告警--", appname}
	// titileString := strings.Join(title, "")
	// msg := []string{"IP地址: ", ip, "<br>主机名称：", hostname, "<br>", "DI.平台服务检查进程存活信息失败：", appname, "<br>告警时间：", date, " ", time}
	// msgString := strings.Join(msg, "")
	// receivers := receiverJson.Receivers

	// wmsg.Title = titileString
	// wmsg.Content = msgString
	// wmsg.Receivers = receivers

	// 解析要发送的消息结构体struct为json格式
	// 将结构体转换为JSON字节流
	log.Println("[I] wmsg's strcut is : ", wmsg)
	log.Printf("[I] wmsg's strcut is : %+v\n", wmsg)
	jsonBytes, err := json.Marshal(wmsg)
	log.Println("[I] jsonBytes 001--> ", string(jsonBytes))
	if err != nil {
		log.Println("[E] WaingMsgSend Error marshaling to JSON:", err)
	}

	return jsonBytes
}

func TemplateJsonWriter(jsonBytes []byte) {
	// 文件写入 -----------------------------------------------------------------------------------------------------------
	filename := "template.json"
	log.Println("filename is :", filename)
	filePath := "../outerlibs/" + filename
	log.Println("[I] outerlibs template filepath is : ", filePath)

	// 使用os.Stat()检查文件是否存在
	_, errosStat := os.Stat(filePath)

	if errosStat != nil {
		// 判断错误是否为文件不存在（os.IsNotExist(err)）
		if os.IsNotExist(errosStat) {
			fmt.Printf("[E]WaingMsgSend File does not exist : %s\n", filePath)

			// 文件不存在，创建一个新的CSV writer用于写入文件
			file, err := os.Create(filePath)
			if err != nil {
				log.Println("WaingMsgSend 文件不存在，创建文件时报错： ", err)
			}
			defer file.Close()

			// 写入数据到文件中
			_, err = file.Write(jsonBytes)
			if err != nil {
				log.Println("[E] WaingMsgSend 文件不存在，创建文件后读取要写入的数据记录报错： ", err)
			}

			log.Printf("IsNotExist json data has been written to %s\n\n", filePath)

		} else {
			fmt.Println("[E] WaingMsgSend Error occurred(else):", errosStat)
		}
	} else {
		log.Printf("[I] WaingMsgSend File exists : %s \n", filePath)

		// 打开文件，如果文件已打开并且是可读的，这通常是允许的
		// file, err := os.Open(filePath)
		// flag := os.O_WRONLY | os.O_CREATE
		flag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		// 设置文件权限，这里是-rwxr-xr-x
		mode := fs.FileMode(0755)

		// 打开文件，即使它已经由其他进程打开用于读取
		file, err := os.OpenFile(filePath, flag, mode)
		// file, err := os.OpenFile(filePath, flag)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		// defer file.Close()

		info, _ := os.Stat(filePath)
		modes := info.Mode()
		if err != nil {
			if os.IsPermission(err) {
				log.Println("Permission denied")
			} else if os.IsExist(err) && (os.FileMode(modes).IsRegular()) { // 检查是否为常规文件且存在
				log.Println("File exists but cannot be opened for some reason") // 可能由于其他原因无法打开，如已被其他进程以独占方式打开等
			} else if os.IsNotExist(err) {
				log.Println("File does not exist")
			} else {
				log.Println("Unknown error:", err)
			}
			return
		}

		defer file.Close()

		// 写入数据到文件中
		log.Println("[W] jsonBytes is :", jsonBytes)
		log.Println("[W] jsonBytes string is :", string(jsonBytes))
		_, err = file.Write(jsonBytes)
		if err != nil {
			log.Println("[E] WaingMsgSend 文件存在，写入数据记录报错： ", err)
		}

		log.Printf("[I] Exist email json data has been written to %s", filePath)
	}
}

func StartEmailSending() {
	// 启动外挂邮件发送包 完成发送邮件
	log.Println("[I] 外挂发送邮件程序")
	dir, err := os.Getwd()
	if err != nil {
		log.Println("[E] Error getting working directory:", err)
		// log.Println("[I]")
		// return
	}
	log.Println("[I]  curr dir is : ", dir)
	// /export/server/monitorXG/bin
	// jarPath := dir + "../outerlibs/uwcdemo-0.0.1-SNAPSHOT.jar"
	// cmd := exec.Command("source", "/etc/profile;", "java", "-jar", "../outerlibs/uwcdemo-0.0.1-SNAPSHOT.jar")
	// cmd := exec.Command("java", "-jar", jarPath)
	// javaCmd.Dir = dir + "/../outerlibs"
	cmd := exec.Command("java", "-jar", "uwcdemo-0.0.1-SNAPSHOT.jar")
	cmd.Dir = dir + "/../outerlibs"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing Java program:", err)
		return
	}
	log.Println(string(output))
	log.Println("[I] 完成运行 外挂告警邮件发送")
}

func WarningCheckFileCreate(msgstruct DataStruct) {
	// 创建静默文件，告警时间设置5分钟静默
	filename := msgstruct.IP + msgstruct.Appname + msgstruct.Date + "txt"
	filepath := "../data/warning/" + filename

	// 10分钟发一次告警邮件，如果10分钟内没有恢复的话
	time := TimeSilenc{
		Time: "10",
	}
	jsonBytes, err := json.Marshal(time)
	if err != nil {
		log.Println("[E] Error WarningCheckFileCreate marshaling to JSON:", err)
		// return
	}
	// 打开或创建json文件以写入
	file, err := os.Create(filepath)
	if err != nil {
		log.Println("[E] Error WarningCheckFileCreate creating file:", err)
		// return
	}
	defer file.Close()

	// 将json字节流写入文件
	_, err = file.Write(jsonBytes)
	if err != nil {
		log.Println("Error writing to file:", err)
		// return
	}
}

func WarningCheckFileDelete(msgstruct DataStruct) {
	filename := msgstruct.IP + msgstruct.Appname + msgstruct.Date + "txt"
	filepath := "../data/warning/" + filename

	// 读取文件，将文件解析到struct
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Println("[E] WarningCheckFileDelete read file get error： ", err)
	}

	defer jsonFile.Close()

	// jsonBytes, err := ioutil.ReadAll(jsonFile)
	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println("[E] WarningCheckFileDelete ioutil.ReadAll(jsonFile) get error： ", err)
	}

	// 关闭文件，防止删除时报错
	jsonFile.Close()

	var timer TimeSilenc

	e := json.Unmarshal(jsonBytes, &timer)
	if err != nil {
		log.Println("[E] json.Unmarshal(jsonBytes, &timer) get error: ", e)
	}

	sleepTime := timer.Time

	sleepTimeInt, err := strconv.ParseInt(sleepTime, 0, 10)

	// 获取文件的创建时间
	// go-fs-util.fs.
	finfo, _ := os.Stat(filepath)
	linuxFileAttr := finfo.Sys().(*syscall.Stat_t) //因为系统运行在linux系统上，使用此方法
	ctime := SecondToTime(linuxFileAttr.Ctim.Sec)

	// 文件创建时间
	// ctimeTime := ctime.Format("16:00")
	// ctimeTimeAfterTenMi := ctime.Add(time.Duration(sleepTimeInt) * time.Minute)
	// 文件创建时间的十分钟后的时间
	// ctimeTimeAfterTenMitime := ctimeTimeAfterTenMi.Format("16:00")

	// 消息中的时间
	msgTime := msgstruct.DateTime
	msgDate := msgstruct.Date

	// 时间字符串
	// timeStr := msgDate + " " + msgTime + " CST"
	timeStr := msgDate + " " + msgTime
	log.Println("[Debug] timeStr value is :", timeStr)

	// 布局字符串
	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai")

	msgTimetime, errTimer := time.ParseInLocation(layout, timeStr, loc)
	if errTimer != nil {
		log.Println("[E] errTimer: ", errTimer)
	}

	// 相差时间
	// timeSub := ctime.Sub(msgTimetime)
	timeSub := msgTimetime.Sub(ctime)
	minutesSub := timeSub.Minutes()

	log.Println("[Debug] minutesSub value is :", minutesSub)

	if minutesSub >= float64(sleepTimeInt) {
		log.Println("[I] 创建时间超过10分钟，删除告警文件文件下一次还是错误的就会创建文件发送告警邮件，进入下一次告警发送")
		errorOfRemove := os.Remove(filepath)
		if errorOfRemove != nil {
			log.Println("[E] 检查文件超过静默时间，删除告警文件时报错： ", errorOfRemove)
		}
	} else {
		log.Println("[E] 未超过静默时间，继续静默")
		return
	}

	// linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
	// ctime := SecondToTime(linuxFileAttr.Ctim.Sec)
	// log.Println("文件创建时间: ", ctime)
	//fmt.Println("最后访问时间", SecondToTime(linuxFileAttr.Atim.Sec))
	//fmt.Println("最后修改时间", SecondToTime(linuxFileAttr.Mtim.Sec))
}

// 时间转换
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func (*JsonContentForEmailSending) WaingMsgSend(msgstruct DataStruct) {

	// 先检查文件，检查文件存在静默后 return 退出执行告警发送 -------------------------------------------------------------------
	log.Println("[I] 检查静默文件......")
	checkfileresult := WarningCheck(msgstruct)

	if checkfileresult == false {
		// 告警信息处理
		jsonBytes := WarningETL(msgstruct)
		// 文件写入
		TemplateJsonWriter(jsonBytes)
		// 启动邮件发送程序
		StartEmailSending()
		// 创建静默文件
		WarningCheckFileCreate(msgstruct)
	} else {
		// true --> 文件存在，执行静默
		// 静默处理需要进一步晚上，同一个服务请求是每30s就会过来一次，连续多次检查都是进程不存在的情况的话，文件会持续更新吸入，增加了磁盘IO
		log.Println("[I] 检查结果为true,执行静默处理后删除静默文件")
		WarningCheckFileDelete(msgstruct)
	}

}

// 告警恢复邮件发送
func (*JsonContentForEmailSending) RecoverMsgSend(msgstrcut DataStruct) {
	// 先检查文件，检查文件存在静默后 return 退出执行告警发送 -------------------------------------------------------------------
	log.Println("[I] 检查静默文件......")
	checkfileresult := WarningCheck(msgstrcut)

	if checkfileresult == false {
		log.Println("[I] 告警恢复检查告警文件不存在，不执行告警恢复操作.")
	} else {
		// true --> 文件存在，执行静默
		// 静默处理需要进一步晚上，同一个服务请求是每30s就会过来一次，连续多次检查都是进程不存在的情况的话，文件会持续更新吸入，增加了磁盘IO
		log.Println("[I] 告警恢复检查告警文件存在，执行告警恢复操作.")
		RcoverCheckFileDelete(msgstrcut)
	}

}

// 发送告警邮件，删除告警文件
func RcoverCheckFileDelete(msgstrcut DataStruct) {
	log.Println("[I] 开始执行告警恢复邮件操作")

	// 发送恢复邮件
	RecoverChecker(msgstrcut)

	// 删除告警文件
	RecoverFileDel(msgstrcut)

}

// 发送邮件恢复
func RecoverChecker(msgstrcut DataStruct) {
	// 告警信息处理
	jsonBytes := WarningETL(msgstrcut)
	// 文件写入
	TemplateJsonWriter(jsonBytes)
	// 启动邮件发送程序
	StartEmailSending()
}

func RecoverFileDel(msgstrcut DataStruct) {
	filename := msgstrcut.IP + msgstrcut.Appname + msgstrcut.Date + "txt"
	filepath := "../data/warning/" + filename
	errOfRemoveFile := os.Remove(filepath)
	if errOfRemoveFile != nil {
		log.Println("[I] RecoverFileDel Error removing file:", errOfRemoveFile)
		// return
	}
}
