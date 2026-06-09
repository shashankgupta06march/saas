USE chatbot_saas;

CREATE TABLE IF NOT EXISTS lead_capture_config (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chatbot_id BIGINT NOT NULL UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    title VARCHAR(255) NOT NULL DEFAULT 'Before we begin...',
    subtitle VARCHAR(500) NOT NULL DEFAULT 'Please share a few details so we can assist you better.',
    fields JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS leads (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chatbot_id BIGINT NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    field_values JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE,
    INDEX idx_leads_chatbot_id (chatbot_id),
    INDEX idx_leads_session_id (session_id),
    INDEX idx_leads_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
