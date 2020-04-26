配置中心（apollo）Go语言客户端使用手册

1. 引入库 `import "github.com/bingoohuang/typhon4g"`
1. 在init函数中/或者main开始增加注册配置文件的监听事件，例如:

   ```go
   func init() {
       // 注册viper读取
       ty := typhon4g.LoadStart()
       ty.Register("application.properties", &typhon4g.ViperListener{})
   }
   ```
   
1. 增加配置文件 `$PWD/etc/typhon-context.properties`

    ```properties
    # 应用程序id
    appId = a100
    # 配置中心服务器类型: apollo, typhon, 默认typhon
    serverType = apollo
    # Meta服务器地址列表，多个时以英文逗号分隔
    metaServers = http://127.0.0.1:11683
    # 配置服务器列表，多个时以英文逗号分隔（如果meta不填，则需要填写此项）
    configServers = 
    ```

1. 使用配置

    ```go
    fmt.Println("name:", viper.GetString("name"))
    fmt.Println("age:", viper.GetInt("age"))
    ```
