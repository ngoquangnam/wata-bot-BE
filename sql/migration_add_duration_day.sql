-- Migration: Add duration_day column to user_bot_subscription table
-- This migration adds a new column to store the selected duration day from the API subscribe request

-- Add duration_day column as VARCHAR
ALTER TABLE `user_bot_subscription` 
ADD COLUMN `duration_day` VARCHAR(20) NOT NULL DEFAULT '0' COMMENT 'Selected duration day from API' AFTER `duration_days`;

-- Update existing records: if duration_days is a JSON array, try to extract the first value as string
-- Otherwise, set to '0' (will need manual update if needed)
UPDATE `user_bot_subscription` 
SET `duration_day` = COALESCE(CAST(JSON_EXTRACT(`duration_days`, '$[0]') AS CHAR), '0')
WHERE `duration_day` = '0' OR `duration_day` = '';

