package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

func CreateMqttTask(mqttData MqttT, c chan MqttMessageT, wg *sync.WaitGroup) error {
	defer wg.Done()
	isRun := true
	server := "tcp://" + mqttData.Host + ":" + mqttData.Port
	opts := mqtt.NewClientOptions().AddBroker(server)
	if mqttData.IsLogin {
		opts.SetUsername(mqttData.Username).SetPassword(mqttData.Password)
	}
	opts.SetClientID(mqttData.ClientId).SetCleanSession(mqttData.IsCleanSession)
	opts.SetKeepAlive(300 * time.Second).SetPingTimeout(30 * time.Second)
	opts.SetAutoReconnect(true).SetConnectRetry(true)
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("MQTT Connected => " + server)
		token := client.Subscribe(mqttData.SubTopic, byte(mqttData.SubQos), func(client mqtt.Client, msg mqtt.Message) {
			//log.Println(server+" 收到订阅消息:", msg.Topic(), string(msg.Payload()))
			c <- MqttMessageT{
				Mqtt:    mqttData,
				Client:  client,
				Message: msg,
			}
		})
		if token.Wait() && token.Error() != nil {
			log.Println("MQTT Subscribe Fatal:", token.Error())
			return
		}
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("Connect loss: %v\n", err)
		isRun = false
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("MQTT Connect Fatal:", token.Error())
		return token.Error()
	}
	for isRun {
		time.Sleep(time.Second)
	}
	log.Println("Close MQTT Task:", mqttData.Host+":"+mqttData.Port)
	return nil
}
