package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"os"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

const topic = "Dorm/Hum"

func main() {
	var broker = "120.78.88.87"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("Dorm/Hum")
	opts.SetPassword("Dorm/Hum")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	//publish(client)
	for {
		publish(client)
	}

}

func publish(client mqtt.Client) {
	num := 10
	min := 75
	max := 99
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("%d.%d", rand.Intn(max-min)+min, rand.Intn(9))
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second * 3)
	}
}

func sub(client mqtt.Client) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func receiveMessages(client mqtt.Client) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Print("Subscribe topic " + topic + " success\n")
}
