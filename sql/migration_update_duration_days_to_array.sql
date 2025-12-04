-- Migration: Update duration_days to JSON array format
-- This migration updates the duration_days column in both bot and user_bot_subscription tables
-- to store a JSON array [5, 15, 30, 60, 90, 180] instead of a single integer value

-- First, alter the column type to JSON in bot table
ALTER TABLE `bot` 
MODIFY COLUMN `duration_days` JSON NOT NULL COMMENT 'Duration in days array';

-- Update existing records in bot table to have the default array [5, 15, 30, 60, 90, 180]
UPDATE `bot` 
SET `duration_days` = JSON_ARRAY(5, 15, 30, 60, 90, 180)
WHERE `duration_days` IS NULL OR JSON_TYPE(`duration_days`) != 'ARRAY';

-- Alter the column type to JSON in user_bot_subscription table
ALTER TABLE `user_bot_subscription` 
MODIFY COLUMN `duration_days` JSON NOT NULL COMMENT 'Duration in days array';

-- Update existing records in user_bot_subscription table to have the default array [5, 15, 30, 60, 90, 180]
-- If duration_days is an integer, convert it to array format
UPDATE `user_bot_subscription` 
SET `duration_days` = JSON_ARRAY(5, 15, 30, 60, 90, 180)
WHERE `duration_days` IS NULL OR JSON_TYPE(`duration_days`) != 'ARRAY';

