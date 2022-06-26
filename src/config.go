package main

import (
	"gopkg.in/ini.v1"
	"log"
)

type ConfigT struct {
	DbType     string
	DbHost     string
	DbPort     int
	DbName     string
	DbUsername string
	DbPassword string
	DbPath     string
	MqttList   []MqttT
}

type MqttT struct {
	Table          string
	Host           string
	Port           string
	ClientId       string
	IsCleanSession bool
	IsLogin        bool
	Username       string
	Password       string
	SubTopic       string
	SubQos         int
}

func LoadConfig(path string) (ConfigT, error) {
	if isOK, err := PathExists(path); isOK == false {
		if err == nil {
			log.Panicln("Config file path is not exist.\n", path)
		} else {
			log.Panicln(err.Error(), path)
		}
	}
	log.Println("Load Config File Path: ", path)
	cfg, err := ini.Load(path)
	if err != nil {
		log.Panicln("Fail to read file: ", err)
	}
	cfg.BlockMode = false //只读操作，增加读取性能
	config := &ConfigT{
		DbType: "sqlite",
		DbPath: "./database.sqlite",
	}
	//log.Println(cfg.Section("main").KeysHash())
	if err = cfg.Section("").MapTo(config); err != nil {
		return *config, err
	}
	for i, mqttCfg := range cfg.Sections() {
		if i == 0 {
			continue
		}
		mqttData := &MqttT{
			Table:          mqttCfg.Name(),
			ClientId:       RandChar(32),
			IsCleanSession: true,
			IsLogin:        false,
			SubTopic:       "#",
			SubQos:         0,
		}
		err := cfg.Section(mqttCfg.Name()).MapTo(mqttData)
		if err != nil {
			log.Println("mqttCfg.MapTo Fatal:", err.Error())
			continue
		}
		config.MqttList = append(config.MqttList, *mqttData)
	}
	return *config, err
}
