package main

import (
	"os"
	"log"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/thedevsaddam/gojsonq"
	"time"
	"bufio"
)

// 获取百度 token
func baidutoken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	grant_type := "client_credentials" // 必须参数，固定 client_credentials
	client_id := "xxxxxxx"  //必须参数，应用的API Key
	client_secret := "xxxxxx"  //必须参数，应用的Secret Key
	resp, err := http.Post(url,"application/x-www-form-urlencoded",strings.NewReader("grant_type="+grant_type+"&client_id="+client_id+"&client_secret="+client_secret))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	res :=gojsonq.New().FromString(string(body)).Find("access_token")  // token
	return res.(string)
 }

 // 图片 base64
 func GetByteFromFile(filePath string)(escapeUrl string, ok bool){
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
	return escapeUrl ,true
}


// general_basic通用文字识别
func general_basic(access_token,image,filePath string){
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
	resp, err := http.Post(url,"application/x-www-form-urlencoded",strings.NewReader("access_token="+access_token+"&image="+image))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	jq :=gojsonq.New().FromString(string(body)).From("words_result").Select("words")
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
		WriteWithIoutil(deviceInfoMap["words"].(string)+"\n",filePath)
	}
	// return string(body)
 }

// accurate_basic通用文字识别高精度版
func accurate_basic(access_token,image,filePath string){
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
	resp, err := http.Post(url,"application/x-www-form-urlencoded",strings.NewReader("access_token="+access_token+"&image="+image))
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	jq :=gojsonq.New().FromString(string(body)).From("words_result").Select("words")
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
		WriteWithIoutil(deviceInfoMap["words"].(string)+"\n",filePath)
	}
	// return string(body)
 }


// 写入文件
func WriteWithIoutil(data,filePath string) {
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
	// 获取百度 token
	access_token := baidutoken()
    // 读取图片目录中的所有文件
    files, err := ioutil.ReadDir("图片/")
    if err != nil {
        panic(err)
	}
	filePath := "识别结果.txt"
	ioutil.WriteFile(filePath,[]byte(""), 0666)
	// 获取文件，并输出它们的名字
    for _, file := range files {
		fmt.Println("------------------------------\n图片： ",file.Name(),"\n内容如下：\n")
		WriteWithIoutil("------------------------------\n图片： "+file.Name()+"\n内容如下：\n",filePath)
		image,err := GetByteFromFile("图片/"+file.Name())
		if !err{
			fmt.Println(err)
		}
		general_basic(access_token,image,filePath)    // 通用文字识别
		// accurate_basic(access_token,image,filePath)  // 通用文字识别高精度版
		time.Sleep(500000000)  // QPS 限制  2次/秒

    }
}
