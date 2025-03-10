CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

DROP TABLE IF EXISTS pastes;
CREATE TABLE pastes (
    id VARCHAR(10) PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255) DEFAULT 'Untitled',
    language VARCHAR(50) DEFAULT 'text',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    views INT DEFAULT 0,
    visibility ENUM('public', 'unlisted') DEFAULT 'public',
    status ENUM('active', 'expired') DEFAULT 'active'
); 