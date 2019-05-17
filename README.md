# typhon4g
typhon client for golang

# Usage

1. Setup `$PWD/etc/typhon-context.properties` like the following:
 
    ```properties
    # 应用程序id
    appID = a100
    
    # Meta服务器地址列表，多个时以英文逗号分隔，默认http://127.0.0.1:11683
    metaServers = http://127.0.0.1:11683
    # 配置刷新时间间隔（秒），默认300
    configRefreshIntervalSeconds = 300
    # 配置刷新是读取超时时间（毫秒），默认5000
    # configReadTimeoutMillis = 5000
    # Http连接超时时间（毫秒），默认1000
    # connectTimeoutMillis = 1000
    # 长连接保持时间（毫秒），默认70000，服务器端是一分钟 （60000），这里配置比服务器端多一点点即可
    # pollingReadTimeoutMillis = 70000
    
    # 快照存储目录，不配置时，默认 ~/.typhon-client/snapshots
    # snapshotsDir = ~/.typhon-client/snapshots
    snapshotsDir = ./etc/snapshots
    
    # meta server/config server连接不上时，等待重试的时间，默认60秒
    # retryNetworkSleepSeconds = 60
    
    # 配置主动刷新时间间隔，默认300秒 （5分钟）
    # configRefreshIntervalSeconds = 300
    
    # Meta主动刷新时间间隔，默认300秒（5分钟）
    # metaRefreshIntervalSeconds = 300
    ```

1. Use the api to access config

    ```go
	prop, err := typhon4g.GetProperties("hello.properties")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("name:", prop.String("name"))
	fmt.Println("home:", prop.StringDefault("home", "中国"))
	fmt.Println("age:", prop.Int("age"))
	fmt.Println("adult", prop.Bool("adult"))
 
 	hello, _ := typhon4g.GetConfFile("hello.json")
 	fmt.Println("hello.json:", hello.Raw())
    ```
    
2. Add the listener to the change of config file

    ```go
    type MyListener struct{}
    
    // Make sure that MyListener implements the interface typhon4g.ConfFileChangeListener
    var _ typhon4g.ConfFileChangeListener = (*MyListener)(nil)
    
    func (l MyListener) OnChange(event typhon4g.ConfFileChangeEvent) (msg string, ok bool) {
        fmt.Println("OnChange", event)
        // eat your own dog food here
        return "your message", true /*  true to means changed OK */
    }
 
 
    // In your code, register the listener instance
    var listener MyListener
    prop.Register(&listener)
 
    ```