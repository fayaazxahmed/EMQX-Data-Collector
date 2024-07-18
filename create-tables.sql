DROP TABLE IF EXISTS emqx_messages;

CREATE TABLE emqx_messages(
    sensor_id         INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    sensor_name VARCHAR(128) NOT NULL,
    topic_name VARCHAR(128) NOT NULL,
    sensor_location VARCHAR(255) NOT NULL,
    last_measured VARCHAR(128) NOT NULL,
    measurement_interval DECIMAL(5,2) NOT NULL
);
