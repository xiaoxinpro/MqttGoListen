# MqttGoListen

MqttGoListen是一款基于Golang开发的MQTT订阅监听并写入数据库的工具，支持 MQTT 3.1 和 MQTT 3.1.1，支持Mysql、Postgre、Sqlite、Memory数据库。

## 开发状态

目前MqttGoListen处于开发初期，主要应用于测试环境。

因此配置文件可能因版本不同而不兼容，请阅读指定版本的配置文档。

* [master](https://github.com/xiaoxinpro/MqttGoListen/tree/master) 分支属于开发版本（稳定性没有经过验证，不提供二进制文件）
* [release](https://github.com/xiaoxinpro/MqttGoListen/releases) 版本为经测试稳定的版本（建议下载最新的 release 版本部署）

Docker版本将待功能基本稳定后进行开发

## 部署说明

MqttGoListen 采用 **应用程序 + 配置文件** 的方式启动，因此配置文件再部署过程中是不可或缺的，关于配置文件的详情请参考[配置文件](#配置文件)章节。

### 二进制文件部署

首先在 [Release](https://github.com/xiaoxinpro/MqttGoListen/releases) 页面下载最新二进制文件，需要根据部署系统选择相应的架构二进制文件。

然后在项目创建一个 `config.ini` 文件

```ini
DbType="sqlite"
DbPath="./db.sqlite"
[test_mqtt]
Host="iot.eclipse.org"
Port=1883
```

最后启动监听程序

```bash
./MqttGoListen -c config.ini
```

## 配置文件

目前配置文件只采用`ini`的格式，其中

* **默认分区** 为全局配置项，如数据库配置等。
* **自定义分区** 为MQTT配置项，分区名自定义但不能重复。

### 全局配置

全局配置必须放在**默认分区**，否则将被忽略执行。

|      键      |        值         | 备注                                |
|:-----------:|:----------------:|-----------------------------------|
|   DbType    |     "mysql"      | 数据库类型（sqlite,mysql,postgre,memory） |
|   DbHost    |   "localhost"    | 数据库主机，仅mysql,postgre需要此项          |
|   DbPort    |       3306       | 数据库端口号，仅mysql,postgre需要此项         |
|   DbName    | "mqtt_go_listen" | 数据库名称，仅mysql,postgre需要此项          |
| DbUsername  |     "admin"      | 数据库用户名，仅mysql,postgre需要此项         |
| DbPassword  |     "123456"     | 数据库密码，仅mysql,postgre需要此项          |
|   DbPath    |  "./db.sqlite"   | 数据库文件路径，仅sqlite需要此项               |

### MQTT配置
|      键      |        值         | 备注                        |
|:-----------:|:----------------:|---------------------------|
| Table          | string | 数据库表名，默认为配置分区名称           |
| Host           | string | MQTT主机地址，必填               |
| Port           | string | MQTT端口号，默认为`1883`           |
| ClientId       | string | MQTT客户端ID，默认为随机字符串        |
| IsCleanSession | bool | MQTT会话清除，默认为`true`          |
| IsLogin        | bool | MQTT用户名密码登录，默认为`false`      |
| Username       | string | MQTT用户名，只有`IsLogin`为`true`时有效 |
| Password       | string | MQTT密码，只有`IsLogin`为`true`时有效  |
| SubTopic       | string | 订阅主题，默认为#                 |
| SubQos         | int | 订阅主题的Qos，取值范围（0，1，2），默认为`0` |

