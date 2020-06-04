# goorc
百度文字识别 API  go语言版

# 流程
获取 access_token
读取相对目录 `图片/` 下的所有文件， 进行 base64 编码 urlencode编码 ，post方式调用 API接口
返回的值 提取存到  `识别结果.txt` 文件下


# 使用

将源码 `goorc.go` 中 `client_id`、`client_secret` 改为自己的 百度文字识别应用
```
	client_id := "xxxxxxx"  //必须参数，应用的API Key
	client_secret := "xxxxxx"  //必须参数，应用的Secret Key
```

支持通用文字识别和通用文字识别高精度版  自行修改
```
		general_basic(access_token,image,filePath)    // 通用文字识别
		// accurate_basic(access_token,image,filePath)  // 通用文字识别高精度版
```

# 运行

```
go run goorc.go
```

# 构建 exe

```
go build goorc.go
```
