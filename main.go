package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const filename = "data.csv"
var cstZone = time.FixedZone("UTC", 8*3600) // 东八
var token string

type dingMsg struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		content string
	} `json:"text"`
	At struct {
		isAtAll int
	} `json:"at"`
}

func log(message interface{}) {
	pc, file, line, _ := runtime.Caller(1)

	f := runtime.FuncForPC(pc).Name()
	now := time.Now().In(cstZone).Format("15:04:05")
	date := time.Now().In(cstZone).Format("2018-10-09")

	str := fmt.Sprintf("%s %s:%d [%s]: %v", now, file, line, f, message)

	fmt.Println(str)

	fname := fmt.Sprintf("ding_eat_%s.log", date)
	logfile, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer logfile.Close()
	if err != nil {
		fmt.Printf("%s %s:%d [%s]: [日志创建错误] %v\r\n", now, file, line, f, err)
	}
	logfile.WriteString(str + "\r\n")

}

func get_weekday(name string) string {
	var weekday string
	switch name {
	case "Monday":
		weekday = "周一"
	case "Tuesday":
		weekday = "周二"
	case "Wednesday":
		weekday = "周三"
	case "Thursday":
		weekday = "周四"
	case "Friday":
		weekday = "周五"
	case "Saturday":
		weekday = "周六"
	case "Sunday":
		weekday = "周日"
	}
	return weekday
}

func get_message(name string) (string, string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log(fmt.Sprintf("文件解析失败：%s", err))
	}
	text := csv.NewReader(strings.NewReader(string(file)))
	lines, _ := text.ReadAll()
	var weekday, msg string
	for i := 0; i < len(lines); i++ {
		if lines[i][0] == name {
			if len(lines[i]) == 2 {
				weekday, msg = lines[i][0], lines[i][1]
			} else {
				log(fmt.Sprintf("配置文件有误"))
				os.Exit(1)
			}
		}
	}
	return weekday, msg
}

func make_msg(weekday, msg string) string {
	var message = ""
	hour, min, _ := time.Now().In(cstZone).Clock()
	now_time := "上午"
	now_min_sec := fmt.Sprintf("%d:%d", hour, min)
	if hour > 12 {
		now_time = "下午"
	}
	message = fmt.Sprintf("现在是 %s %s %s啦, %s", weekday, now_time, now_min_sec, msg)
	log(fmt.Sprintf("消息已生成:%s", message))
	return message

}

func send_msg(msg string) {
	var dingmsg dingMsg
	dingmsg.At.isAtAll = 1
	dingmsg.MsgType = "text"
	dingmsg.Text.content = msg
	JsonMsg, err := json.Marshal(dingmsg)
	if err != nil {
		log(err)
		os.Exit(1)
	}
	body := bytes.NewBuffer([]byte(JsonMsg))
	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log(err)
		os.Exit(1)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log(err)
		os.Exit(1)
		return
	}
	log(fmt.Sprintf("接口返回信息：%s", result))
	return
}

func init() {
	flag.StringVar(&token, "t", "", "token when u create a robot in dingtalk")
	flag.Parse()
	log(fmt.Sprintf("初始化完成，获取token为：%s", token))
}

func main() {
	send_msg(make_msg(get_message(get_weekday(time.Now().In(cstZone).Weekday().String()))))
}
