package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// 获取百度 token
func baidutoken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	grant_type := "client_credentials"                  // 必须参数，固定 client_credentials
	client_id := "V0dawQ7U4b3uQt8FSwlrD1bS"             //必须参数，应用的API Key
	client_secret := "qSYxzgTB4DSB08IHXQt209MhdVuyrf3I" //必须参数，应用的Secret Key
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("grant_type="+grant_type+"&client_id="+client_id+"&client_secret="+client_secret))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return jsoniter.Get(body, "access_token").ToString()
	//fmt.Println(string(body))
	//res :=gojsonq.New().From(string(body)).Find("access_token")  // token
	//return res.(string)
}

// 图片 base64
func GetByteFromFile(filePath string) (escapeUrl string, ok bool) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	pic, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	encodeToString := base64.StdEncoding.EncodeToString(pic)
	escapeUrl = url.QueryEscape(encodeToString)
	return escapeUrl, true
}

// general_basic通用文字识别
func general_basic(filename,access_token, image string) (string, error) {
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("access_token="+access_token+"&image="+image))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("============")
	//fmt.Println(string(body))
	//fmt.Println("============")

	words:=string(body)

	if strings.Contains(words, "1D:") {
		theIndex := strings.Index(words, "1D:")
		id := words[theIndex+3 :theIndex+19]
		return id,nil
	}

	return "",errors.New(filename+"找不到 ID")



	//jq :=gojsonq.New().From(string(body)).From("words_result").Select("words")
	//deviceInfoList, ok := jq.Get().([]interface{})
	//if !ok {
	//	fmt.Println("Convert deviceInfoList error")
	//}
	//for _, deviceInfo := range deviceInfoList {
	//	deviceInfoMap, ok := deviceInfo.(map[string]interface{})
	//	if !ok {
	//		fmt.Println("Convert deviceInfoMap error")
	//	}
	//	fmt.Println(deviceInfoMap["words"])
	//	WriteWithIoutil(deviceInfoMap["words"].(string)+"\n",filePath)
	//}
	// return string(body)
}

// accurate_basic通用文字识别高精度版
func accurate_basic(access_token, image, filePath string) {
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("access_token="+access_token+"&image="+image))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	jq := gojsonq.New().From(string(body)).From("words_result").Select("words")
	deviceInfoList, ok := jq.Get().([]interface{})
	if !ok {
		fmt.Println("Convert deviceInfoList error")
	}
	for _, deviceInfo := range deviceInfoList {
		deviceInfoMap, ok := deviceInfo.(map[string]interface{})
		if !ok {
			fmt.Println("Convert deviceInfoMap error")
		}
		fmt.Println(deviceInfoMap["words"])
		WriteWithIoutil(deviceInfoMap["words"].(string)+"\n", filePath)
	}
	// return string(body)
}

// 写入文件
func WriteWithIoutil(data, filePath string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func main() {

	logFile, err := os.Create("d:\\go.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger := log.New(logFile, "", log.Llongfile)

	logger.Println("开始请求百度 token")

	// 获取百度 token
	access_token := baidutoken()
	logger.Println("百度 token "+access_token)


	// 读取图片目录中的所有文件
	//file, err := ioutil.ReadFile("/Users/jinsong/work/个人信息图/mmexport1611754750303.jpg")
	//if err != nil {
	//    panic(err)
	//}

	//dir := "/Users/jinsong/work/information"
	dir := "d:\\information"

	files,err :=ioutil.ReadDir(dir)
	if err != nil {
		logger.Fatalln("无法读取目录")
		panic(err)
	}

	logger.Println("目录 "+dir)


	for _, file := range files {





		if strings.Contains(file.Name(), ".") {
			theIndex := strings.Index(file.Name(), ".")
			ext := file.Name()[theIndex :len(file.Name())]
			target := dir+"\\"+file.Name()
			image,err1 := GetByteFromFile(target)
			if !err1{
				logger.Fatalln("无法读取文件",err1)
				continue
			}
			id,err2 :=general_basic(target,access_token,image)    // 通用文字识别
			// accurate_basic(access_token,image,filePath)  // 通用文字识别高精度版


			logger.Println("target "+target)
			logger.Println("ext "+ext)


			if err2 != nil {
				logger.Fatalln("无法读取ID",err2)
				continue
			}

			logger.Println("id "+id)

			body, err3 := ioutil.ReadFile(target)
			if err3 != nil {
				logger.Fatalln("无法写入文件",err3)
				continue
			}

			dest := dir+"1\\"+id+ext

			logger.Println("dest "+dest)

			ioutil.WriteFile(dest,body,0666)
			fmt.Println(target +" -> "+dest)
			time.Sleep(500000000)  // QPS 限制  2次/秒


		}















	}


	//
	//image, err := GetByteFromFile("/Users/jinsong/work/个人信息图/mmexport1611754750303.jpg")
	//if !err {
	//	fmt.Println(err)
	//}
	//id := general_basic(access_token, image) // 通用文字识别

	//fmt.Println(words)
	//
	//if strings.Contains(words, "1D:") {
	//	theIndex := strings.Index(words, "1D:")
	//	id := words[theIndex+3 :theIndex+19]
	//
	//	fmt.Println(id)
	//}

	os.Exit(0)
	//filePath := "识别结果.txt"
	//ioutil.WriteFile(filePath,[]byte(""), 0666)
	//// 获取文件，并输出它们的名字
	//for _, file := range files {
	//	fmt.Println("------------------------------\n图片： ",file.Name(),"\n内容如下：\n")
	//	WriteWithIoutil("------------------------------\n图片： "+file.Name()+"\n内容如下：\n",filePath)
	//	image,err := GetByteFromFile("图片/"+file.Name())
	//	if !err{
	//		fmt.Println(err)
	//	}
	//	general_basic(access_token,image,filePath)    // 通用文字识别
	//	// accurate_basic(access_token,image,filePath)  // 通用文字识别高精度版
	//	time.Sleep(500000000)  // QPS 限制  2次/秒
	//
	//}
}


//截取字符串 start 起点下标 end 终点下标(不包括)
func substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
