-- ========================================
-- TẠO DATABASE & SỬ DỤNG
-- ========================================
CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

-- ========================================
-- XÓA CÁC BẢNG NẾU TỒN TẠI
-- ========================================
DROP TABLE IF EXISTS paste_views;
DROP TABLE IF EXISTS paste;

-- ========================================
-- BẢNG paste: chứa nội dung paste
-- ========================================
CREATE TABLE paste (
  id VARCHAR(10) NOT NULL,
  content TEXT NOT NULL,
  title VARCHAR(255) DEFAULT 'Untitled',
  language VARCHAR(50) DEFAULT 'text',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP DEFAULT NULL,
  visibility ENUM('PUBLIC','UNLISTED') DEFAULT 'PUBLIC',
  PRIMARY KEY (id)
);

-- ========================================
-- BẢNG paste_views: log lượt xem từng paste
-- ========================================
CREATE TABLE paste_views (
  id INT NOT NULL AUTO_INCREMENT,
  paste_id VARCHAR(10) NOT NULL,
  viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_viewed_at_paste_id (viewed_at, paste_id),
  FOREIGN KEY (paste_id) REFERENCES paste(id) ON DELETE CASCADE
);

-- ========================================
-- USER CREATION AND PRIVILEGES
-- ========================================

-- DROP existing users to recreate with correct passwords if needed
DROP USER IF EXISTS 'monitor'@'%';
DROP USER IF EXISTS 'proxysql_user'@'%';

-- Create monitor user with the SAME password as in ProxySQL config
CREATE USER 'monitor'@'%' IDENTIFIED BY 'StrongSecurePassword123';

-- Grant monitoring privileges to monitor user
GRANT SELECT, REPLICATION CLIENT, PROCESS ON *.* TO 'monitor'@'%';

-- Create application user for ProxySQL
CREATE USER 'proxysql_user'@'%' IDENTIFIED BY 'proxysql_password';

-- Grant privileges to proxysql_user
GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'proxysql_user'@'%';
GRANT ALL PRIVILEGES ON pastebin.* TO 'proxysql_user'@'%';

-- Apply all privilege changes
FLUSH PRIVILEGES;