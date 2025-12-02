-- Create database
CREATE DATABASE IF NOT EXISTS `wata_bot` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `wata_bot`;

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

-- Create bot table
CREATE TABLE IF NOT EXISTS `bot` (
  `id` VARCHAR(20) NOT NULL COMMENT 'Bot ID',
  `name` VARCHAR(100) NOT NULL COMMENT 'Bot name',
  `icon_letter` VARCHAR(1) NOT NULL COMMENT 'Icon letter',
  `risk_level` VARCHAR(50) NOT NULL COMMENT 'Risk level',
  `duration_days` INT NOT NULL COMMENT 'Duration in days',
  `expected_return_percent` INT NOT NULL COMMENT 'Expected return percentage',
  `apr_display` VARCHAR(100) NOT NULL COMMENT 'APR display text',
  `min_investment` INT NOT NULL COMMENT 'Minimum investment',
  `max_investment` INT NOT NULL COMMENT 'Maximum investment',
  `investment_range` VARCHAR(50) NOT NULL COMMENT 'Investment range',
  `subscribers` INT NOT NULL DEFAULT 0 COMMENT 'Number of subscribers',
  `author` VARCHAR(100) NOT NULL COMMENT 'Author name',
  `description` TEXT NOT NULL COMMENT 'Bot description',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Is bot active',
  `lockup_period` VARCHAR(50) NOT NULL COMMENT 'Lockup period',
  `expected_return` VARCHAR(50) NOT NULL COMMENT 'Expected return',
  `min_investment_display` VARCHAR(50) NOT NULL COMMENT 'Min investment display',
  `max_investment_display` VARCHAR(50) NOT NULL COMMENT 'Max investment display',
  `roi30d` VARCHAR(50) NOT NULL COMMENT 'ROI 30 days',
  `win_rate` VARCHAR(50) NOT NULL COMMENT 'Win rate',
  `trading_pair` VARCHAR(200) NOT NULL COMMENT 'Trading pair',
  `total_trades` INT NOT NULL DEFAULT 0 COMMENT 'Total trades',
  `pnl30d` DECIMAL(20, 2) NOT NULL DEFAULT 0.00 COMMENT 'P&L 30 days',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  PRIMARY KEY (`id`),
  KEY `idx_is_active` (`is_active`),
  KEY `idx_risk_level` (`risk_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot table';

-- Create user_bot_subscription table
CREATE TABLE IF NOT EXISTS `user_bot_subscription` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Subscription ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
  `bot_id` VARCHAR(20) NOT NULL COMMENT 'Bot ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_bot` (`user_id`, `bot_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_bot_id` (`bot_id`),
  CONSTRAINT `fk_subscription_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_subscription_bot` FOREIGN KEY (`bot_id`) REFERENCES `bot` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User bot subscription table';

