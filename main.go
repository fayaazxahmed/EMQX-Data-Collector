package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	//"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-sql-driver/mysql"
)

// Struct containing the message information that is inserted into the database. Primary key and last measured timestamp are updated automatically so they are not included
type Message struct {
	topic_name  string
	measurement string
	sensor_name string
}

var db *sql.DB

// Update constants for specific deployment being connected to
const protocol = "ssl"
const broker = "g332f11e.ala.eu-central-1.emqxsl.com"
const port = 8883
const topic = "root/faux/data/#"
const username = "Fayaaz"
const password = "FA5"

func main() {
	client := createMqttClient()
	time.Sleep(time.Second * 10)                                  // pause minimum of 2 seconds to wait for the subscription function to be ready, otherwise subscriber function doesn't receive messages
	var broker_msg, broker_topic, sensor_name = subscribe(client) // we use goroutine to run the subscription function, and store the returned message data in various variables

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"), //Set DBUSER and DBPASS environment variables
		Passwd:               os.Getenv("DBPASS"), //Alternatively, the SQL username and password can also be set manually without using environment variables but that will make them visible to the public if published to a public repository
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "emqx_data",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	msg := Message{broker_topic, broker_msg, sensor_name}
	tableInsert(db, msg)

}

func createMqttClient() mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	clientID := "emqx_cloude096fd"

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()

	opts.AddBroker(connectAddress)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(10*time.Second) && token.Error() != nil {
		log.Printf("\nConnection error: %s\n", token.Error())
	}
	return client
}

func subscribe(client mqtt.Client) (string, string, string) {
	qos := 0
	broker_msg := make(chan string)
	broker_topic := make(chan string)
	sensor_name := make(chan string)
	client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s, from topic: %s \n", msg.Payload(), msg.Topic())
		split := strings.Split(string(msg.Topic()), "/")
		broker_msg <- string(msg.Payload())
		broker_topic <- string(msg.Topic())
		sensor_name <- split[3]
	})

	return (<-broker_msg), (<-broker_topic), (<-sensor_name)
}

/*
func tableInsert(db *sql.DB, message Message) int {
	query := `INSERT INTO emqx_messages (topic_name, measurement, last_measured)
		VALUES (?, ?, NOW())`

	last_entry, err := db.Exec(query, message.topic_name, message.measurement)
	if err != nil {
		log.Fatal(err)
	}

	//Gets the ID for the last entry in the table
	lastInsertId, err := last_entry.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(lastInsertId)
}
*/

func tableInsert(db *sql.DB, message Message) int {
	topicIDQuery := `SELECT topicID FROM Topics WHERE topicName = ?`
	var topicID int
	//Check if that topic is already inside the Topics table and retrieves it's topicID if it is in the table
	err := db.QueryRow(topicIDQuery, message.topic_name).Scan(&topicID)
	if err != nil {
		//Inserts topic into the table if it's a new topic
		query := `INSERT INTO Topics (topicName)
				VALUES (?)`

		lastTopicID, err := db.Exec(query, message.topic_name)
		if err != nil {
			log.Fatal(err)
		}

		//Gets the ID for the last entry in the table
		lastInsertId, err := lastTopicID.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New topic added with ID: %d\n", lastInsertId)
	}

	//After checking the Topics table to see if the topic exists, insert the message data and timestamp into the data, as well as the topic ID as a foreign key
	logInsertQuery := `INSERT INTO Logs (topicID, measurement, measureTime) VALUES (?, ?, NOW())`
	logInsert, err := db.Exec(logInsertQuery, topicID, message.measurement)
	if err != nil {
		log.Fatal(err)
	}
	lastInsertID, err := logInsert.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(lastInsertID)
}
