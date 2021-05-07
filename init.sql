CREATE TABLE IF NOT EXISTS C_APP (
	USER_ID varchar(40) NOT NULL,
	RECORD_ID varchar(40) NOT NULL,
	NAME varchar(70) NOT NULL,
	DESCRIPTION varchar(1000) NOT NULL,
	STATUS varchar(50) NOT NULL,
	CREATE_TIME DATETIME NOT NULL,
	CANCEL_TIME DATETIME,
	PRIMARY KEY (RECORD_ID)
);

CREATE TABLE IF NOT EXISTS C_APP_TOPIC (
	RECORD_ID varchar(40) NOT NULL,
	APP_ID varchar(40) NOT NULL,
	NAME varchar(100) NOT NULL,
	DESCRIPTION varchar(100) NOT NULL,
	PRIMARY KEY (RECORD_ID)
);

CREATE TABLE IF NOT EXISTS C_APP_USER (
	RECORD_ID varchar(40) NOT NULL,
	CLIENT_ID varchar(40) NOT NULL,
	APP_ID varchar(40) NOT NULL,
	STATUS varchar(50) NOT NULL,
	REGISTER_TIME DATETIME NOT NULL,
	LAST_STATUS_CHANGE_TIME DATETIME NOT NULL,
	PRIMARY KEY (RECORD_ID,CLIENT_ID)
);

CREATE TABLE IF NOT EXISTS C_APP_TOPIC_USER (
	RECORD_ID varchar(40) NOT NULL,
	APP_ID varchar(40) NOT NULL,
	TOPIC_ID varchar(40) NOT NULL,
	`USER_ID` varchar(40) NOT NULL,
	PRIMARY KEY (RECORD_ID)
);

ALTER TABLE C_APP_TOPIC ADD CONSTRAINT C_APP_TOPIC_fk0 FOREIGN KEY (APP_ID) REFERENCES C_APP(RECORD_ID);

ALTER TABLE C_APP_USER ADD CONSTRAINT C_APP_USER_fk0 FOREIGN KEY (APP_ID) REFERENCES C_APP(RECORD_ID);

ALTER TABLE C_APP_TOPIC_USER ADD CONSTRAINT C_APP_TOPIC_USER_fk0 FOREIGN KEY (TOPIC_ID) REFERENCES C_APP_TOPIC(RECORD_ID);

ALTER TABLE C_APP_TOPIC_USER ADD CONSTRAINT C_APP_TOPIC_USER_fk1 FOREIGN KEY (USER_ID) REFERENCES C_APP_USER(RECORD_ID);

ALTER TABLE C_APP_TOPIC_USER ADD CONSTRAINT C_APP_TOPIC_USER_fk2 FOREIGN KEY (APP_ID) REFERENCES C_APP(RECORD_ID);
