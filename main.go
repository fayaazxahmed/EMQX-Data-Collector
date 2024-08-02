package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	//"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const protocol = "ssl"
const broker = "g332f11e.ala.eu-central-1.emqxsl.com"
const port = 8883
const topic = "root/faux/data/#"
const username = "Fayaaz"
const password = "FA5"

func main() {
	client := createMqttClient()
	go subscribe(client)        // we use goroutine to run the subscription function
	time.Sleep(time.Second * 1) // pause 1s to wait for the subscription function to be ready
	subscribe(client)
}

func createMqttClient() mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	//rand.Seed(time.Now().UnixNano())
	clientID := "emqx_cloude096fd" //fmt.Sprintf("go-client-%d", rand.Int())

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()

	opts.AddBroker(connectAddress)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)

	opts.SetTLSConfig(loadTLSConfig("emqxsl-ca.pem"))
	opts.SetTLSConfig(loadTLSConfig("main.go"))

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(10*time.Second) && token.Error() != nil {
		log.Printf("\nConnection error: %s\n", token.Error())
	}
	return client
}

func subscribe(client mqtt.Client) {
	qos := 0
	client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s ,from topic: %s \n", msg.Payload(), msg.Topic())
	})
}

func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}
