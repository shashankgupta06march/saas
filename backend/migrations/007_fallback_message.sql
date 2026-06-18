-- Admin-configurable message shown to the user when the chatbot finds nothing
-- relevant in its knowledge base. When empty, the bot answers as before.
ALTER TABLE chatbot_settings
  ADD COLUMN fallback_message VARCHAR(1000) NULL AFTER icon_url;
