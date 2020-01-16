# 通用透传下位机服务系统

用于部署搭建透传方式的物联网下位机。已实现以下功能：

- [x] 以TCP/IP透传方式与传感器通讯，包括发送指令和接收传感器数据
- [x] MQTT方式与上位机通讯
- [x] 通过conf.json文件, 本地管理传感器设备
- [x] 通过conf.json文件, 远程管理传感器设备
- [x] 下位机Level日志上报
- [x] 传感器状态上报
- [x] 传感器异常机制
- [x] 传感器开关
- [x] 下位机重启
- [x] 具有缓存功能，在与上位机通讯失败时，缓存未推送成功的数据。与传感器通讯失败时，缓存失败的指令
- [ ] 提供Web界面配置下位机（同上）
- [ ] 提供Web界面管理传感器设备（使用远程Web界面设置）
- [ ] 与上位机通讯(暂时以http协议通讯)，把传感器数据推送给上位机服务器，接收上位机指令（使用MQTT通讯）

##### 部署方式

1. 使用支持MQTT协议的消息中间件, 这里以RabbitMQ为例进行部署
```cmd
docker run -d \
--name some-rabbit \
-p 5672:5672 -p 15672:15672 \
-e RABBITMQ_DEFAULT_USER=user \
-e RABBITMQ_DEFAULT_PASS=password \
rabbitmq:3-management

# 进入bash开启mqtt插件
rabbitmq-plugins enable rabbitmq_mqtt
```

2. 修改下位机指向的消息中间件, 即 `cnf/conf.json` 文件下的地址/端口/协议等参数, 格式如下:
```json
{
  # 中间件地址
  "broker_ip": "106.13.79.157",
  # 中间件端口
  "broker_port": "1883",
  # 中间件协议
  "broker_scheme": "tcp",
  # 中间件用户名
  "broker_username": "r3inb",
  # 中间件密码
  "broker_password": "159463",
  # 可缺省
  # "broker_client_id": "",
}
```

3. 完成传感器匹配, 修改 `cnf/conf.json` 文件, 完成sensor的配置, 格式如下:
```json
{
  # 下位机名称
  "name": "示例收集器名",
  # MQTT IP
  # MQTT PORT
  # MQTT SCHME
  "localSensorInformation": [
    {
      # 传感器物理地址
      "addr": 6,
      # 传感器类型, 类型表在sensor_query.go下
      "type": 0,
      # 传感器依附的DTU IP
      "attach": "172.20.10.4",
      # 测量间隔时间
      "interval": 10,
      # 传感器ID
      "sensorID": "7eb220dd-6127-58c7-8663-bf2f55371b78"
    }
  ]
}
```

4. 启动程序/docker compose

#### 表格

##### mqtt_client

| 函数名 | 描述                    |  返回值 |
| ------------- | ------------------------------ |----------------------|
| `GetMQTTInstance()`| 获得一个MQTT的连接 | mqtt.Client, error |


##### LocalSensorInformation

| 成员变量      | 描述 |
| --------- | -----:|
| TaskHandler     |   自定义传感器任务 |
| Status | 状态 |
| Addr  | 传感器设备地址 |
| Type     |   传感器类型 |
| Attach      |    传感器附着的透传设备 |
| Interval  | 最大间隔时间(秒) |
| SensorID     |   传感器ID |

	
| 函数名 | 描述                    |  返回值 |
| ---------------------- | ------------------------------ |----------------------|
| `CreateTask(times int)`| 创建一个传感器任务 | error |
| `RemoveTask()`| 移除一个传感器任务 | error |
| `UpdateTask(times int)`| 更新一个传感器任务 | error |
| `AddTaskHandler(callback Job)`| 创建自定义任务 |  |
| `RemoveTaskHandler()`| 移除自定义任务 | bool |
| `Open()`| 开启传感器 |  |
| `Close()`| 关闭传感器 |  |


##### LocalDeviceList

| 成员变量      | 描述 |
| --------- | -----:|
| Name     |   透传设备名称 |
| LocalSensorInformation      |    传感器集合 |

| 函数名 | 描述                    |  返回值 |
| ---------------------- | ------------------------------ |----------------------|
| `GetLocalSensorList(attachIP string)`| 获取相关传感器数据列表 | []LocalSensorInformation |


##### TaskSensorKey

| 成员变量      | 描述 |
| --------- | -----:|
| Addr     |   设备地址 |
| Attach  | 附着设备 |
| Type     |   指令类型 |


##### TaskSensorBody

| 成员变量      | 描述 |
| --------- | -----:|
| TaskSensorKey     |   任务唯一id |
| Type  | 指令类型 |
| RequestData     |   生成的指令数据 |
| SensorID     |   传感器ID |

| 函数名 | 描述                    |  返回值 |
| ---------------------- | ------------------------------ |----------------------|
| `CreateMeasureRequest()`| 创建测量请求数据 |  |


#### 使用

使用 `MQTTMapping(topic string, callback mqtt.MessageHandler)` 进行主题订阅, 如下:
在新增服务时, 新增 `mqtt.MessageHandler` 即可
```go
    MQTTMapping("sensor/action/clear", sensor.ClearExceptionHandler)
```

使用 `MQTTPublish(topic string, payload interface{})` 进行主题发布

文档建设中...

备注: 
l. 临时MQ后台 http://106.13.79.157:15672/ 账号/密码: admin (有效期至2020.02.06)
