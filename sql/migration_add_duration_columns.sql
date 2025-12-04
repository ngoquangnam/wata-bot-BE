-- Migration: Add duration_days and duration_day columns to user_bot_subscription table
-- This migration checks if columns exist and adds them if they don't

-- Check and add duration_days column (JSON) if it doesn't exist
SET @col_exists = 0;
SELECT COUNT(*) INTO @col_exists 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'user_bot_subscription' 
  AND COLUMN_NAME = 'duration_days';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `user_bot_subscription` ADD COLUMN `duration_days` JSON NOT NULL DEFAULT (JSON_ARRAY(5, 15, 30, 60, 90, 180)) COMMENT ''Duration in days array'' AFTER `bot_id`',
    'SELECT ''Column duration_days already exists'' AS message');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Check and add duration_day column (VARCHAR) if it doesn't exist
SET @col_exists = 0;
SELECT COUNT(*) INTO @col_exists 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'user_bot_subscription' 
  AND COLUMN_NAME = 'duration_day';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `user_bot_subscription` ADD COLUMN `duration_day` VARCHAR(20) NOT NULL DEFAULT ''0'' COMMENT ''Selected duration day from API'' AFTER `duration_days`',
    'SELECT ''Column duration_day already exists'' AS message');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Update existing records: set default values if needed
UPDATE `user_bot_subscription` 
SET `duration_days` = JSON_ARRAY(5, 15, 30, 60, 90, 180)
WHERE `duration_days` IS NULL OR JSON_TYPE(`duration_days`) != 'ARRAY';

UPDATE `user_bot_subscription` 
SET `duration_day` = '0'
WHERE `duration_day` IS NULL OR `duration_day` = '';

