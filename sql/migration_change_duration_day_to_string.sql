-- Migration: Change duration_day column from INT to VARCHAR
-- This migration changes the duration_day column type from INT to VARCHAR(20)

-- If column exists as INT, alter it to VARCHAR
ALTER TABLE `user_bot_subscription` 
MODIFY COLUMN `duration_day` VARCHAR(20) NOT NULL DEFAULT '0' COMMENT 'Selected duration day from API';

-- Update existing records: convert integer values to string
UPDATE `user_bot_subscription` 
SET `duration_day` = CAST(`duration_day` AS CHAR)
WHERE `duration_day` IS NOT NULL;

