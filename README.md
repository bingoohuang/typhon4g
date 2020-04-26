# typhon4g

[![Travis CI](https://img.shields.io/travis/bingoohuang/typhon4g/master.svg?style=flat-square)](https://travis-ci.com/bingoohuang/typhon4g)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/bingoohuang/typhon4g/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/bingoohuang/typhon4g)
[![Coverage Status](http://codecov.io/github/bingoohuang/typhon4g/coverage.svg?branch=master)](http://codecov.io/github/bingoohuang/typhon4g?branch=master)
[![goreport](https://www.goreportcard.com/badge/github.com/bingoohuang/typhon4g)](https://www.goreportcard.com/report/github.com/bingoohuang/typhon4g)

typhon/apollo client for golang

## Setup config file like the following:
 
```properties
# $PWD/etc/typhon-context.properties

# 应用程序id
appId = a100

# 配置中心服务器类型: apollo, typhon, 默认typhon
serverType = apollo


# 以下是apollo专有
cluster    = default
dataCenter = 
# localIp不配置的话，默认获取en0/eth0网卡的v4版本的IP地址
localIp    =

# Meta服务器地址列表，多个时以英文逗号分隔
metaServers = http://127.0.0.1:11683
# 配置服务器列表，多个时以英文逗号分隔（如果meta不填，则需要填写此项）
configServers = 
# 配置刷新时间间隔，默认5分钟
configRefreshInterval = 5m
# 配置刷新是读取超时时间，默认5s
# configReadTimeout = 5s
# Http连接超时时间，默认60s
# connectTimeout = 60s
# 长连接保持时间，默认70秒，服务器端是一分钟，这里配置比服务器端多一点点即可
# pollingReadTimeout = 70s

# 快照存储目录，不配置时，默认 ~/.typhon-client/snapshots
# snapshotsDir = ~/.typhon-client/snapshots
snapshotsDir = ./etc/snapshots

# meta server/config server连接不上时，等待重试的时间，默认60秒
# retryNetworkSleep = 60s

# 配置主动刷新时间间隔，默认5分钟
# configRefreshInterval = 5m

# Meta主动刷新时间间隔，默认5分钟
# metaRefreshInterval = 5m

# 推送配置所需要的用户名密码
# postAuth = admin:admin

# 服务端证书校验用的根证书文件(不校验服务端证书时，请注释本配置项)
#rootPem = ./root.pem

# 提供给服务端校验的客户端证书文件和客户端私钥文件（服务端不需要校验时，请注释两个配置项）
#clientPem = ./client.pem
#clientKey = ./client.key
```

## Use the api to access config

```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
	"github.com/bingoohuang/typhon4g"
)
 
func init() {
    // 注册viper读取
    ty := typhon4g.LoadStart()
	ty.Register("application.properties", &typhon4g.ViperListener{Prefix:""})
	ty.Register("hello.yaml", &typhon4g.ViperListener{Prefix:"hello."})
}

func main() {
    fmt.Println("name:", viper.GetString("name"))
    fmt.Println("age:", viper.GetInt("hello.age"))
    fmt.Println("adult", viper.GetBool("hello.adult"))
}
```
