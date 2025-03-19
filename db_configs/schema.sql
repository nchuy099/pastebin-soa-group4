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
    status ENUM('active', 'expired') DEFAULT 'active',

    INDEX idx_created_at (created_at),      -- Speeds up queries by creation date
    INDEX idx_expires_at (expires_at),      -- Helps in expiration-based queries
    INDEX idx_status (status),              -- Speeds up filtering by status
    INDEX idx_visibility (visibility),      -- Optimizes visibility-based lookups
    INDEX idx_language (language)           -- Speeds up language-based searches
); 