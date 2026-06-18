-- Add a launcher icon image URL for the chatbot widget button.
-- This is separate from avatar_url (used for the in-message bot avatar).
ALTER TABLE chatbot_settings
  ADD COLUMN icon_url VARCHAR(500) NULL AFTER avatar_url;
