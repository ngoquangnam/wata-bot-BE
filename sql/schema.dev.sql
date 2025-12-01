-- Create database for development
CREATE DATABASE IF NOT EXISTS `wata_bot_dev` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `wata_bot_dev`;

-- Create user table
CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `address` VARCHAR(42) NOT NULL COMMENT 'Wallet address',
  `referral_code` VARCHAR(8) NOT NULL COMMENT 'Referral code',
  `invite_code` VARCHAR(42) DEFAULT NULL COMMENT 'Invite code used',
  `wata_reward` INT NOT NULL DEFAULT 0 COMMENT 'WATA reward points',
  `role` VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT 'User role',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_address` (`address`),
  KEY `idx_referral_code` (`referral_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User table';

