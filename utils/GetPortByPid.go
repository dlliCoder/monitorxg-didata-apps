package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

func GetPortsByPID(pid int) ([]string, error) {
	// Linux环境下从/proc读取端口信息
	filename := fmt.Sprintf("/proc/%d/net/tcp", pid)
	tcpBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	filename = fmt.Sprintf("/proc/%d/net/udp", pid)
	udpBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ports []string

	// 简化的示例：假设每行都是有效的，并只提取本地端口
	for _, content := range [][]byte{tcpBytes, udpBytes} {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines[:len(lines)-1] {
			// 去掉最后一行（通常是空行）
			parts := strings.Fields(line)
			log.Println("parts is : ", parts)
			if len(parts) > 3 && parts[3] == "0A" {
				// 判断是否为监听状态（状态列0A）
				_, portStr, _ := net.SplitHostPort(parts[1])
				log.Println("net.SplitHostPort(parts[1]) --> ", portStr)
				decimalNumber, err := strconv.ParseInt(portStr, 16, 0)
				if err != nil {
					log.Println("[E] decimalNumber get error")
				}
				// port, _ := strconv.Atoi(portStr)
				log.Println("net.SplitHostPort(parts[1]) --> strconv.Atoi(portStr)--> ", decimalNumber)
				ports = append(ports, strconv.Itoa(int(decimalNumber)))
				log.Println("net.SplitHostPort(parts[1]) --> strconv.Atoi(portStr)--> append(ports, strconv.Itoa(port))--> ", ports)
			}
		}
	}

	return ports, nil
}
