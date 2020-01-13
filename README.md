# 通用透传下位机服务系统

用于部署搭建透传方式的物联网下位机。主要有以下功能：

- [x] 以TCP/IP透传方式与传感器通讯，包括发送指令和接收传感器数据
- [x] 通过conf.json文件, 管理传感器设备，绑定设备的IP地址
- [x] 分配传感器设备ID
- [ ] 传感器状态上报 -> coming soon
- [ ] 远程设置传感器 -> coming soon
- [ ] 提供Web界面管理传感器设备
- [ ] 提供Web界面配置下位机
- [ ] 与上位机通讯(暂时以http协议通讯)，把传感器数据推送给上位机服务器，接收上位机指令
- [x] 具有缓存功能，在与上位机通讯失败时，缓存未推送成功的数据。与传感器通讯失败时，缓存失败的指令

##### 使用支持MQTT协议的消息中间件, 这里以RabbitMQ为例
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

#####mqtt_client
| 函数名 | 描述                    |  返回值 |
| ------------- | ------------------------------ |----------------------|
| `GetMQTTInstance()`| 获得一个MQTT的连接 | mqtt.Client, error |


#####LocalSensorInformation

| 成员变量      | 描述 |
| --------- | -----:|
| TaskHandler     |   自定义传感器任务 |
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


#####LocalDeviceList
| 成员变量      | 描述 |
| --------- | -----:|
| Name     |   透传设备名称 |
| ID  | 透传设备地址 |
| IP     |   透传设备IP |
| LocalSensorInformation      |    传感器集合 |

| 函数名 | 描述                    |  返回值 |
| ---------------------- | ------------------------------ |----------------------|
| `GetLocalSensorList(attachIP string)`| 获取相关传感器数据列表 | []LocalSensorInformation |


#####TaskSensorKey
| 成员变量      | 描述 |
| --------- | -----:|
| Addr     |   设备地址 |
| Attach  | 附着设备 |
| Type     |   指令类型 |


#####TaskSensorBody
| 成员变量      | 描述 |
| --------- | -----:|
| TaskSensorKey     |   任务唯一id |
| Type  | 指令类型 |
| RequestData     |   生成的指令数据 |
| SensorID     |   传感器ID |

| 函数名 | 描述                    |  返回值 |
| ---------------------- | ------------------------------ |----------------------|
| `CreateMeasureRequest()`| 创建测量请求数据 |  |

文档建设中...

备注: 
l. 临时MQ后台 http://106.13.79.157:15672/ 账号/密码: admin (有效期至2020.02.06)
2. 重构
3. 读写协程
4. 超时反馈 -> 传感器的状态标识
5. 任务队列与连接的整合


