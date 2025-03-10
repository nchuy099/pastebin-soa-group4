CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

CREATE TABLE IF NOT EXISTS pastes (
    id VARCHAR(10) PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255) DEFAULT 'Untitled',
    language VARCHAR(50) DEFAULT 'text',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    views INT DEFAULT 0,
    visibility ENUM('public', 'private', 'unlisted') DEFAULT 'public'
); 