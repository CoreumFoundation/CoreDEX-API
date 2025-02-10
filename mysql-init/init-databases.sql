CREATE DATABASE IF NOT EXISTS friendly_dex;
CREATE USER IF NOT EXISTS 'testuser'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON friendly_dex.* TO 'testuser'@'%';
GRANT SUPER ON *.* TO 'testuser'@'%';
FLUSH PRIVILEGES;

USE friendly_dex;

CREATE TABLE IF NOT EXISTS State (
		StateType INT, 
		Content TEXT, 
		MetaData JSON, 
		Network INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))),
		UNIQUE KEY (Network,StateType)
	);


INSERT INTO State (StateType, Content, MetaData) VALUES (1, '{"Height":12696406}', '{"Network": 3, "UpdatedAt": {"seconds": 1738799304, "nanos": 164479000}}');
