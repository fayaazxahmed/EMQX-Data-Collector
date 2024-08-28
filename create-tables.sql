/*
DROP TABLE IF EXISTS emqx_messages;
CREATE TABLE emqx_messages (
  sensor_id INT AUTO_INCREMENT PRIMARY KEY,
  topic_name VARCHAR(128) NOT NULL,
  measurement VARCHAR(128) NOT NULL,
  last_measured timestamp DEFAULT NOW()
);
*/

DROP TABLE IF EXISTS Topics;
CREATE TABLE Topics (
  topicID INT AUTO_INCREMENT,
  -- sensorID INT,
  topicName VARCHAR(128) NOT NULL,
  -- gasName VARCHAR(128),
  -- unitOfMeasurement VARCHAR(128) NOT NULL,
  PRIMARY KEY (topicID)
  -- FOREIGN KEY (sensorID) REFERENCES Sensors(sensorID)
);

DROP TABLE IF EXISTS Logs;
CREATE TABLE Logs (
  logID INT AUTO_INCREMENT,
  topicID INT,
  measurement VARCHAR(128) NOT NULL,
  measureTime timestamp DEFAULT NOW(),
  PRIMARY KEY (logID),
  FOREIGN KEY (topicID) REFERENCES Topics(topicID)
);