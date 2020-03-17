# typhon4g
typhon/apollo client for golang [![Go Report Card](https://goreportcard.com/badge/github.com/bingoohuang/typhon4g)](https://goreportcard.com/report/github.com/bingoohuang/typhon4g)


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
localIp    =

# Meta服务器地址列表，多个时以英文逗号分隔
metaServers = http://127.0.0.1:11683
# 配置服务器列表（如果meta不填，则需要填写此项）
configServers = 
# 配置刷新时间间隔，默认5分钟
configRefreshInterval = 5m
# 配置刷新是读取超时时间，默认5s
# configReadTimeout = 5s
# Http连接超时时间，默认1秒
# connectTimeout = 1s
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
 
var ty = typhon4g.LoadStart()

func init() {
    // 注册viper读取
	ty.Register("application.properties", &typhon4g.ViperListener{})
}

func main() {
    fmt.Println("name:", viper.GetString("name"))
    fmt.Println("age:", viper.GetInt("age"))
    fmt.Println("adult", viper.GetBool("adult"))
}
```
