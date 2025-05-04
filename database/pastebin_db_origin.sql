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
  PRIMARY KEY (id),
  INDEX idx_expires_at (expires_at),
  INDEX idx_expires_at_id (expires_at, id)
);

-- ========================================
-- BẢNG paste_views: log lượt xem từng paste
-- ========================================
CREATE TABLE paste_views (
  id INT NOT NULL AUTO_INCREMENT,
  paste_id VARCHAR(10) NOT NULL,
  viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_viewed_at_paste_id (viewed_at, paste_id)
);

