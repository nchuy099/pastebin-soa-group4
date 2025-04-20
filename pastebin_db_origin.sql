CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

DROP TABLE IF EXISTS paste;

CREATE TABLE paste (
  id varchar(10) NOT NULL,
  content text NOT NULL,
  title varchar(255) DEFAULT 'Untitled',
  language varchar(50) DEFAULT 'text',
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  expires_at datetime DEFAULT NULL,
  views int DEFAULT 0,
  visibility enum('PUBLIC','UNLISTED') DEFAULT 'PUBLIC',
  PRIMARY KEY (id),
  INDEX idx_created_at (created_at),
  INDEX idx_id_expires (id, expires_at)

) ENGINE=InnoDB;  