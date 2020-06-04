# goorc
百度文字识别 API  go语言版

# 流程
获取 access_token
读取相对目录 `图片/` 下的所有文件， 进行 base64 编码 urlencode编码 ，post方式调用 API接口
返回的值 提取存到  `识别结果.txt` 文件下

# 运行

```
go run goorc.go
```

# 构建 exe

```
go build goorc.go
```
