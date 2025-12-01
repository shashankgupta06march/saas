-- Seed Data Migration
-- This file creates demo organizations and users for testing and deployment
-- Password hashes are generated using bcrypt with cost 10
-- Note: Database name should be specified when running this migration

-- Seed Organizations
INSERT IGNORE INTO organizations (name, api_key, plan_type, status, created_at, updated_at) VALUES
('Test Company', '06efecc90c6069f61c803ca415cdb01d5a7fe3c192f10a6bd70031e3bda484e6', 'free', 'active', NOW(), NOW()),
('Demo Corp', 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2', 'premium', 'active', NOW(), NOW()),
('Enterprise Solutions Ltd', 'e1f2g3h4i5j6k7l8m9n0o1p2q3r4s5t6u7v8w9x0y1z2a3b4c5d6e7f8g9h0i1j2', 'enterprise', 'active', NOW(), NOW()),
('Startup Inc', 'f1g2h3i4j5k6l7m8n9o0p1q2r3s4t5u6v7w8x9y0z1a2b3c4d5e6f7g8h9i0j1k2', 'free', 'active', NOW(), NOW()),
('Tech Innovations', 'g1h2i3j4k5l6m7n8o9p0q1r2s3t4u5v6w7x8y9z0a1b2c3d4e5f6g7h8i9j0k1l2', 'premium', 'active', NOW(), NOW()),
('Global Services', 'h1i2j3k4l5m6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1c2d3e4f5g6h7i8j9k0l1m2', 'enterprise', 'active', NOW(), NOW()),
('Small Business Co', 'i1j2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2e3f4g5h6i7j8k9l0m1n2', 'free', 'trial', NOW(), NOW()),
('Medium Enterprise', 'j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m0n1o2', 'premium', 'active', NOW(), NOW());

-- Seed Users with bcrypt hashed passwords
-- Note: These hashes are pre-generated for the passwords shown in comments
-- password123 -> $2a$10$rI4qR5F5vJ5vJ5vJ5vJ5vOqR5F5vJ5vJ5vJ5vJ5vJ5vJ5vJ5vJ5v
-- demo123456 -> $2a$10$demohashdemohashdemohashdemohashdemohashdemohashdemoha
-- etc.

-- For Test Company
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(1, 'admin@test.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW()),
(1, 'user@test.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user', NOW());

-- For Demo Corp
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(2, 'demo@democorp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW()),
(2, 'manager@democorp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'manager', NOW());

-- For Enterprise Solutions Ltd
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(3, 'admin@enterprise.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW()),
(3, 'support@enterprise.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'support', NOW());

-- For Startup Inc
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(4, 'contact@startup.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW());

-- For Tech Innovations
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(5, 'admin@techinnovations.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW());

-- For Global Services
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(6, 'admin@globalservices.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW());

-- For Small Business Co
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(7, 'owner@smallbiz.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW());

-- For Medium Enterprise
INSERT IGNORE INTO users (organization_id, email, password_hash, role, created_at) VALUES
(8, 'admin@mediumenterprise.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW());

-- Seed Sample Chatbots
INSERT IGNORE INTO chatbots (organization_id, name, description, status, created_at, updated_at) VALUES
(1, 'Test Company Support Bot', 'Customer support chatbot for Test Company', 'active', NOW(), NOW()),
(2, 'Demo Corp Support Bot', 'Customer support chatbot for Demo Corp', 'active', NOW(), NOW()),
(3, 'Enterprise Solutions Ltd Support Bot', 'Customer support chatbot for Enterprise Solutions Ltd', 'active', NOW(), NOW());

-- Seed Chatbot Settings
INSERT IGNORE INTO chatbot_settings (chatbot_id, theme_color, position, welcome_message, widget_size) VALUES
(1, '#007bff', 'bottom-right', 'Hi! Welcome to Test Company. How can I help you today?', 'medium'),
(2, '#28a745', 'bottom-right', 'Hi! Welcome to Demo Corp. How can I help you today?', 'medium'),
(3, '#6f42c1', 'bottom-right', 'Hi! Welcome to Enterprise Solutions. How can I help you today?', 'large');

-- Seed Sample Knowledge Base Entries
INSERT IGNORE INTO knowledge_base (organization_id, chatbot_id, title, content, content_type, status, version, chunking_method, created_at) VALUES
(1, 1, 'Company Information', 'Test Company is a leading provider of innovative solutions. We have been in business since 2020 and serve customers worldwide.', 'text', 'active', 1, 'paragraph', NOW()),
(1, 1, 'Product Features', 'Our main product features include: 1) Easy integration 2) Real-time analytics 3) 24/7 support 4) Scalable architecture 5) Security-first design', 'text', 'active', 1, 'paragraph', NOW()),
(2, 2, 'Pricing Plans', 'We offer three pricing plans: Free (basic features), Premium ($29/month), and Enterprise (custom pricing with dedicated support).', 'text', 'active', 1, 'paragraph', NOW()),
(3, 3, 'Support Hours', 'Our enterprise support team is available 24/7 via email, phone, and live chat. Response time: <1 hour for critical issues.', 'text', 'active', 1, 'paragraph', NOW());

-- ============================================================================
-- SEED DATA SUMMARY
-- ============================================================================
-- 
-- Organizations Created: 8
--   - 3 Free tier
--   - 3 Premium tier
--   - 2 Enterprise tier
-- 
-- Users Created: 11
--   All passwords are: "password123" (for simplicity in demo)
-- 
-- Credentials:
-- ============================================================================
-- | Organization              | Email                        | Role    | Plan       |
-- |---------------------------|------------------------------|---------|------------|
-- | Test Company              | admin@test.com               | admin   | free       |
-- | Test Company              | user@test.com                | user    | free       |
-- | Demo Corp                 | demo@democorp.com            | admin   | premium    |
-- | Demo Corp                 | manager@democorp.com         | manager | premium    |
-- | Enterprise Solutions Ltd  | admin@enterprise.com         | admin   | enterprise |
-- | Enterprise Solutions Ltd  | support@enterprise.com       | support | enterprise |
-- | Startup Inc               | contact@startup.com          | admin   | free       |
-- | Tech Innovations          | admin@techinnovations.com    | admin   | premium    |
-- | Global Services           | admin@globalservices.com     | admin   | enterprise |
-- | Small Business Co         | owner@smallbiz.com           | admin   | free       |
-- | Medium Enterprise         | admin@mediumenterprise.com   | admin   | premium    |
-- ============================================================================
-- 
-- Chatbots Created: 3
-- Knowledge Base Entries: 4
-- 
-- ============================================================================

