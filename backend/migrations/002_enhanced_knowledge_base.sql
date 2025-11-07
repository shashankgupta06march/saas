-- Enhanced Knowledge Base Schema
-- Features: Folders, Tags, Versions, Advanced Metadata, Chunking

-- Create folders/categories table
CREATE TABLE IF NOT EXISTS `kb_folders` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `chatbot_id` BIGINT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `parent_id` BIGINT NULL,
  `color` VARCHAR(7) DEFAULT '#3B82F6',
  `icon` VARCHAR(50) DEFAULT 'folder',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`chatbot_id`) REFERENCES `chatbots`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`parent_id`) REFERENCES `kb_folders`(`id`) ON DELETE CASCADE,
  INDEX `idx_chatbot_folders` (`chatbot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create tags table
CREATE TABLE IF NOT EXISTS `kb_tags` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `chatbot_id` BIGINT NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `color` VARCHAR(7) DEFAULT '#10B981',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`chatbot_id`) REFERENCES `chatbots`(`id`) ON DELETE CASCADE,
  UNIQUE KEY `unique_tag_per_chatbot` (`chatbot_id`, `name`),
  INDEX `idx_chatbot_tags` (`chatbot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create knowledge base tags junction table
CREATE TABLE IF NOT EXISTS `kb_entry_tags` (
  `kb_id` BIGINT NOT NULL,
  `tag_id` BIGINT NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`kb_id`, `tag_id`),
  FOREIGN KEY (`kb_id`) REFERENCES `knowledge_base`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`tag_id`) REFERENCES `kb_tags`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create version history table
CREATE TABLE IF NOT EXISTS `kb_versions` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `kb_id` BIGINT NOT NULL,
  `version` INT NOT NULL,
  `title` VARCHAR(500),
  `content` LONGTEXT NOT NULL,
  `content_type` VARCHAR(50),
  `source_url` VARCHAR(1000),
  `changed_by` BIGINT NULL COMMENT 'User who made the change',
  `change_summary` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`kb_id`) REFERENCES `knowledge_base`(`id`) ON DELETE CASCADE,
  INDEX `idx_kb_versions` (`kb_id`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create chunks table for smart chunking
CREATE TABLE IF NOT EXISTS `kb_chunks` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `kb_id` BIGINT NOT NULL,
  `chunk_index` INT NOT NULL,
  `content` TEXT NOT NULL,
  `embedding_vector` JSON NULL,
  `token_count` INT NULL,
  `metadata` JSON NULL COMMENT 'Chunk-specific metadata (page number, section, etc.)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`kb_id`) REFERENCES `knowledge_base`(`id`) ON DELETE CASCADE,
  INDEX `idx_kb_chunks` (`kb_id`, `chunk_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create sync sources table for auto-sync
CREATE TABLE IF NOT EXISTS `kb_sync_sources` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `chatbot_id` BIGINT NOT NULL,
  `source_type` ENUM('google_drive', 'dropbox', 'notion', 'wordpress', 'confluence', 'webhook') NOT NULL,
  `source_identifier` VARCHAR(500) NOT NULL COMMENT 'Folder ID, URL, etc.',
  `auth_token` TEXT NULL COMMENT 'Encrypted auth token',
  `sync_frequency` ENUM('realtime', 'hourly', 'daily', 'weekly') DEFAULT 'daily',
  `last_sync_at` TIMESTAMP NULL,
  `next_sync_at` TIMESTAMP NULL,
  `is_active` BOOLEAN DEFAULT TRUE,
  `sync_settings` JSON NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`chatbot_id`) REFERENCES `chatbots`(`id`) ON DELETE CASCADE,
  INDEX `idx_chatbot_sync_sources` (`chatbot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create knowledge testing/quality table
CREATE TABLE IF NOT EXISTS `kb_quality_tests` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `chatbot_id` BIGINT NOT NULL,
  `test_query` TEXT NOT NULL,
  `expected_content` TEXT NULL,
  `actual_response` TEXT NULL,
  `relevance_score` DECIMAL(3,2) NULL COMMENT 'Score 0-1',
  `passed` BOOLEAN NULL,
  `tested_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`chatbot_id`) REFERENCES `chatbots`(`id`) ON DELETE CASCADE,
  INDEX `idx_chatbot_tests` (`chatbot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create search index table for better search performance
CREATE TABLE IF NOT EXISTS `kb_search_index` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `kb_id` BIGINT NOT NULL,
  `keyword` VARCHAR(255) NOT NULL,
  `frequency` INT DEFAULT 1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`kb_id`) REFERENCES `knowledge_base`(`id`) ON DELETE CASCADE,
  INDEX `idx_keyword_search` (`keyword`),
  INDEX `idx_kb_keywords` (`kb_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
