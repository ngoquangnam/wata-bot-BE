-- Migration: Add duration_days column to user_bot_subscription table
-- Run this if the user_bot_subscription table already exists

USE `wata_bot`;

-- Add duration_days column if it doesn't exist
ALTER TABLE `user_bot_subscription` 
ADD COLUMN IF NOT EXISTS `duration_days` INT NOT NULL DEFAULT 0 COMMENT 'Duration in days' AFTER `bot_id`;

-- Update existing records to use bot's default duration_days
UPDATE `user_bot_subscription` ubs
INNER JOIN `bot` b ON ubs.bot_id = b.id
SET ubs.duration_days = b.duration_days
WHERE ubs.duration_days = 0 OR ubs.duration_days IS NULL;

