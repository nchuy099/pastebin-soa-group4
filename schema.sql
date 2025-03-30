CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

DROP TABLE IF EXISTS paste;
CREATE TABLE paste (
                        id VARCHAR(10) PRIMARY KEY,
                        content TEXT NOT NULL,
                        title VARCHAR(255) DEFAULT 'Untitled',
                        language VARCHAR(50) DEFAULT 'text',
                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                        expires_at DATETIME NULL,
                        views INT DEFAULT 0,
                        visibility ENUM('PUBLIC', 'UNLISTED') DEFAULT 'PUBLIC',
                        status ENUM('ACTIVE', 'EXPIRED') DEFAULT 'ACTIVE'
);