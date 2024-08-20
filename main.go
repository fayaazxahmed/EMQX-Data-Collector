package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	//"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-sql-driver/mysql"
)

type Message struct {
	topic_name  string
	measurement string
}

var db *sql.DB

const protocol = "ssl"
const broker = "g332f11e.ala.eu-central-1.emqxsl.com"
const port = 8883
const topic = "root/faux/data/#"
const username = "Fayaaz"
const password = "FA5"

func main() {
	client := createMqttClient()
	go subscribe(client)         // we use goroutine to run the subscription function
	time.Sleep(time.Second * 10) // pause minimum of 2 seconds to wait for the subscription function to be ready, otherwise subscriber function won't receive messages
	var broker_msg = subscribe(client)

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "games",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable(db)

	msg := Message{broker_msg, topic}
	tableInsert(db, msg)

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

func subscribe(client mqtt.Client) string {
	qos := 0
	broker_msg := make(chan string)
	client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s, from topic: %s \n", msg.Payload(), msg.Topic())
		broker_msg <- string(msg.Payload())
	})
	return <-broker_msg
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

func createTable(db *sql.DB) {
	query := `DROP TABLE IF EXISTS emqx_messages;
		CREATE TABLE emqx_messages (
			sensor_id INT AUTO_INCREMENT PRIMARY KEY,
  			topic_name VARCHAR(128) NOT NULL,
  			measurement VARCHAR(128) NOT NULL,
  			last_measured timestamp DEFAULT NOW()
		)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func tableInsert(db *sql.DB, message Message) int {
	query := `INSERT INTO message (topic_name, measurement)
		VALUES ($1, $2, $3) RETURNING sensor_id`

	var primary_key int
	err := db.QueryRow(query, message.topic_name, message.measurement).Scan(primary_key)
	if err != nil {
		log.Fatal(err)
	}

	return primary_key
}
