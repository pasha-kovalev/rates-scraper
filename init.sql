CREATE TABLE IF NOT EXISTS rates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cur_id INT,
    rate_date DATETIME,
    cur_abbreviation VARCHAR(10),
    cur_scale INT,
    cur_name VARCHAR(100),
    cur_official_rate DECIMAL(10, 4),
    INDEX idx_rate_date (rate_date),
    UNIQUE KEY unique_rate_date_cur_id(rate_date, cur_id)
);
