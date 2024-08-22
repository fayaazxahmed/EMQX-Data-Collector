DROP TABLE IF EXISTS emqx_messages;
CREATE TABLE emqx_messages (
  sensor_id INT AUTO_INCREMENT PRIMARY KEY,
  topic_name VARCHAR(128) NOT NULL,
  measurement VARCHAR(128) NOT NULL,
  last_measured timestamp DEFAULT NOW()
);