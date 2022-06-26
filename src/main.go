package main

import (
	"flag"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

type MqttMessageT struct {
	Mqtt    MqttT
	Client  mqtt.Client
	Message mqtt.Message
}

// ConfigFilePath 配置文件路径
var (
	ConfigFilePath string
)

func init() {
	// 配置错误提示
	mqtt.DEBUG = log.New(os.Stdout, "[mqtt][DBG]", 0)
	mqtt.ERROR = log.New(os.Stdout, "[mqtt][ERR]", 0)
}

func main() {
	// 获取启动参数
	flag.StringVar(&ConfigFilePath, "c", "", "Config File Path.")
	flag.Parse()

	// 加载配置
	config, err := LoadConfig(ConfigFilePath)
	if err != nil {
		log.Panicln("Load Config Error:", err.Error())
	}
	if len(config.MqttList) == 0 {
		log.Panicln("Load Config Error: MQTT host missing from config file.")
	}

	// 连接数据库
	db, err := InitDataBase(config)
	if err != nil {
		log.Panicln("Initial Database Error:", err)
	}

	// 创建数据库写入任务
	mqttChan := make(chan MqttMessageT, 128)
	go MessageReceiveServer(db, mqttChan)

	// 创建监听任务
	wg := sync.WaitGroup{}
	for _, mqttData := range config.MqttList {
		// 检查或创建表
		err := CheckDbTable(db, mqttData.Table)
		if err != nil {
			log.Println("Database Table Error:", err.Error())
			continue
		}
		// 创建MQTT监听任务
		wg.Add(1)
		go CreateMqttTask(mqttData, mqttChan, &wg)
	}
	wg.Wait()
}

func MessageReceiveServer(db *gorm.DB, c chan MqttMessageT) {
	log.Println("Receive message service start...")
	for data := range c {
		log.Println(data.Mqtt.Host+" -> MQTT Recv:", data.Message.Topic(), string(data.Message.Payload()))
		db.Table(data.Mqtt.Table).Create(&MqttModel{
			Topic:     data.Message.Topic(),
			Payload:   string(data.Message.Payload()),
			Qos:       data.Message.Qos(),
			Retained:  data.Message.Retained(),
			MessageID: data.Message.MessageID(),
			CreatedAt: time.Now(),
		})
	}
	log.Println("Receive message service stop.")
}
