package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 360路由器限制，写了个HostFilter，应该兼容360路由器基本型号的后台Host限制
// 不兼容自己改
// chat.openai.com
// 使用请将浏览器F12手动设置一次Host观察请求，复制为CURL CMD 将请求Body的Host替换为书输出内容，注意，不要替换了不该替换的符号
func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前目录:", err)
		return
	}
	hostsFilePath := filepath.Join(currentDir, "hosts")
	hostsFile, err := os.Open(hostsFilePath)
	if err != nil {
		fmt.Println("无法打开hosts文件:", err)
		return
	}
	defer hostsFile.Close()
	outputFile, err := os.Create("host.txt")
	if err != nil {
		fmt.Println("无法创建输出文件:", err)
		return
	}
	defer outputFile.Close()
	regex := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+([\w\.-]+)`)
	scanner := bufio.NewScanner(hostsFile)
	var outputLines []string
	for scanner.Scan() {
		line := scanner.Text()
		matches := regex.FindStringSubmatch(line)
		if len(matches) == 3 {
			domain := matches[2]
			segments := strings.Split(domain, ".")
			if len(segments) <= 5 && len(domain) <= 63 {
				outputLine := matches[1] + "^%^20" + domain
				outputLines = append(outputLines, outputLine)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取hosts文件时出错:", err)
		return
	}
	outputContent := strings.Join(outputLines, "^%^0A")
	_, err = outputFile.WriteString(outputContent)
	if err != nil {
		fmt.Println("写入输出文件时出错:", err)
		return
	}
	fmt.Println("筛选完成，结果已写入host.txt文件。")
}
