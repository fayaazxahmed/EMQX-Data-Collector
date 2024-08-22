# EMQX Data Collector
#### The EMQX Data Collector program is used to connect to an online EMQX broker to retrieve and store the message data sent through the broker in a local SQL database
This program has been designed primarily to collect and store the readings sent from various sensors and store them locally for easy access in the future. Your computer can subscribe to an MQTT topic to retrieve and store all data sent to that topic 
from any device using this program, allowing for quick and simple storage of IoT device data.

## Setup and Run
#### Make sure you have MySQL installed and an active EMQX deployment
Run the following command to log into MySQL, and enter your system password:
```
mysql -u root -p
```

After logging in, run the following commands to set a username/password and to create and use a new database for holding the tables with the EMQX Data, respectively.
```
mysql> CREATE USER 'username'@'host' IDENTIFIED WITH authentication_plugin BY 'password';
mysql> CREATE DATABASE emqx_data;
mysql> USE DATABASE emqx_data;
```
>[!NOTE]
> __Optional:__ Set two environment variables, DBUSER and DBPASS, to store your SQL username and password.\
> This will keep your login info hidden if you are planning on keeping your code in a public repository.

Now you can start entering the connection details for your deployment in *main.go*

Keeping the protocol and ports unchanged, enter the connection address, topic, and username/password for your EMQX instance:
```
const protocol = "ssl"
const broker =
const port = 8883
const topic =
const username =
const password =
```

>[!NOTE]
> If you are not using environment variables for this program, enter your username and password manually inside the following configuration in the main function in *main.go*:
```
cfg := mysql.Config{
  User:                 os.Getenv("DBUSER"),
  Passwd:               os.Getenv("DBPASS"),
  Net:                  "tcp",
  Addr:                 "127.0.0.1:3306",
  DBName:               "emqx_data",
  AllowNativePasswords: true,
}
```

Finally, run the following command in your project directory and this program will be ready to go
```
go run .
```
