-- Migration: Remove duration_days column from user_bot_subscription table
-- This migration removes the duration_days column and keeps only duration_day

-- Check if duration_days column exists and remove it
SET @col_exists = 0;
SELECT COUNT(*) INTO @col_exists 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'user_bot_subscription' 
  AND COLUMN_NAME = 'duration_days';

SET @sql = IF(@col_exists > 0,
    'ALTER TABLE `user_bot_subscription` DROP COLUMN `duration_days`',
    'SELECT ''Column duration_days does not exist'' AS message');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Ensure duration_day column exists
SET @col_exists = 0;
SELECT COUNT(*) INTO @col_exists 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'user_bot_subscription' 
  AND COLUMN_NAME = 'duration_day';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `user_bot_subscription` ADD COLUMN `duration_day` VARCHAR(20) NOT NULL DEFAULT ''0'' COMMENT ''Selected duration day from API'' AFTER `bot_id`',
    'SELECT ''Column duration_day already exists'' AS message');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

