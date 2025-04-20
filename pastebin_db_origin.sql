-- ========================================
-- TẠO DATABASE & SỬ DỤNG
-- ========================================
CREATE DATABASE IF NOT EXISTS pastebin;
USE pastebin;

-- ========================================
-- XÓA CÁC BẢNG NẾU TỒN TẠI
-- ========================================
DROP TABLE IF EXISTS paste_stats;
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
-- BẢNG thống kê theo tháng
-- ========================================
CREATE TABLE paste_stats (
  month_year CHAR(7) PRIMARY KEY, -- e.g., '2025-04'
  total_views INT DEFAULT 0,
  avg_views_per_paste INT DEFAULT 0,
  min_views INT DEFAULT 0,
  max_views INT DEFAULT 0,
  last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ========================================
-- TẠO VIEW read-only cho bảng thống kê
-- ========================================
CREATE OR REPLACE VIEW paste_stats_readonly AS
SELECT *
FROM paste_stats;

-- ========================================
-- EVENT: Tự động cập nhật mỗi giờ
-- ========================================
DELIMITER $$

CREATE EVENT IF NOT EXISTS refresh_paste_stats
ON SCHEDULE EVERY 1 HOUR
DO
BEGIN
  DECLARE monthStart DATE;
  DECLARE monthEnd DATE;
  SET monthStart = DATE_FORMAT(CURDATE(), '%Y-%m-01');
  SET monthEnd = LAST_DAY(monthStart) + INTERVAL 1 DAY;

  REPLACE INTO paste_stats (
    month_year, total_views, avg_views_per_paste, min_views, max_views
  )
  SELECT
    DATE_FORMAT(monthStart, '%Y-%m'),
    COALESCE(SUM(view_count), 0),
    COALESCE(ROUND(AVG(view_count)), 0),
    COALESCE(MIN(view_count), 0),
    COALESCE(MAX(view_count), 0)
  FROM (
    SELECT paste_id, COUNT(*) AS view_count
    FROM paste_views
    WHERE viewed_at >= monthStart AND viewed_at < monthEnd
    GROUP BY paste_id
  ) AS monthly_paste_views;
END$$

DELIMITER ;

-- ========================================
-- BẬT EVENT SCHEDULER
-- ========================================
SET GLOBAL event_scheduler = ON;

-- ========================================
-- CHẠY CẬP NHẬT NGAY LẬP TỨC
-- ========================================
SET @monthStart = DATE_FORMAT(CURDATE(), '%Y-%m-01');
SET @monthEnd = LAST_DAY(@monthStart) + INTERVAL 1 DAY;

REPLACE INTO paste_stats (
  month_year, total_views, avg_views_per_paste, min_views, max_views
)
SELECT
  DATE_FORMAT(@monthStart, '%Y-%m'),
  COALESCE(SUM(view_count), 0),
  COALESCE(ROUND(AVG(view_count)), 0),
  COALESCE(MIN(view_count), 0),
  COALESCE(MAX(view_count), 0)
FROM (
  SELECT paste_id, COUNT(*) AS view_count
  FROM paste_views
  WHERE viewed_at >= @monthStart AND viewed_at < @monthEnd
  GROUP BY paste_id
) AS monthly_paste_views;
