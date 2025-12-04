-- Migration: Fix empty wata_balance and usdt_balance values
-- Run this to update existing records with empty balance values to '0'

USE `wata_bot`;

-- Update existing records to have default values if they are NULL or empty string
UPDATE `user` SET `wata_balance` = '0' WHERE `wata_balance` IS NULL OR `wata_balance` = '' OR TRIM(`wata_balance`) = '';
UPDATE `user` SET `usdt_balance` = '0' WHERE `usdt_balance` IS NULL OR `usdt_balance` = '' OR TRIM(`usdt_balance`) = '';

