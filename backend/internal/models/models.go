package models

import (
	"database/sql"
	"time"
)

type Organization struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	PlanType  string    `json:"plan_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Chatbot struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ChatbotSettings struct {
	ID             int64  `json:"id"`
	ChatbotID      int64  `json:"chatbot_id"`
	ThemeColor     string `json:"theme_color"`
	Position       string `json:"position"`
	WelcomeMessage string `json:"welcome_message"`
	AvatarURL      string `json:"avatar_url"`
	CustomCSS      string `json:"custom_css"`
	WidgetSize     string `json:"widget_size"`
	Suggestions    string `json:"suggestions"` // JSON array of suggestion strings
}

type KnowledgeBase struct {
	ID              int64           `json:"id"`
	OrganizationID  int64           `json:"organization_id"`
	ChatbotID       int64           `json:"chatbot_id"`
	FolderID        sql.NullInt64   `json:"folder_id"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	ContentType     string          `json:"content_type"`
	SourceURL       string          `json:"source_url"`
	Version         int             `json:"version"`
	Status          string          `json:"status"`
	EmbeddingVector string          `json:"embedding_vector"`
	FileSize        sql.NullInt64   `json:"file_size"`
	PageCount       sql.NullInt32   `json:"page_count"`
	WordCount       sql.NullInt32   `json:"word_count"`
	ChunkSize       int             `json:"chunk_size"`
	ChunkOverlap    int             `json:"chunk_overlap"`
	ChunkingMethod  string          `json:"chunking_method"`
	QualityScore    sql.NullFloat64 `json:"quality_score"`
	LastSyncedAt    sql.NullTime    `json:"last_synced_at"`
	SyncEnabled     bool            `json:"sync_enabled"`
	Metadata        sql.NullString  `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
}

type KBFolder struct {
	ID          int64          `json:"id"`
	ChatbotID   int64          `json:"chatbot_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	ParentID    sql.NullInt64  `json:"parent_id"`
	Color       string         `json:"color"`
	Icon        string         `json:"icon"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type KBTag struct {
	ID        int64     `json:"id"`
	ChatbotID int64     `json:"chatbot_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type KBEntryTag struct {
	KBID      int64     `json:"kb_id"`
	TagID     int64     `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

type KBVersion struct {
	ID            int64          `json:"id"`
	KBID          int64          `json:"kb_id"`
	Version       int            `json:"version"`
	Title         sql.NullString `json:"title"`
	Content       string         `json:"content"`
	ContentType   sql.NullString `json:"content_type"`
	SourceURL     sql.NullString `json:"source_url"`
	ChangedBy     sql.NullInt64  `json:"changed_by"`
	ChangeSummary sql.NullString `json:"change_summary"`
	CreatedAt     time.Time      `json:"created_at"`
}

type KBChunk struct {
	ID              int64          `json:"id"`
	KBID            int64          `json:"kb_id"`
	ChunkIndex      int            `json:"chunk_index"`
	Content         string         `json:"content"`
	EmbeddingVector sql.NullString `json:"embedding_vector"`
	TokenCount      sql.NullInt32  `json:"token_count"`
	Metadata        sql.NullString `json:"metadata"`
	CreatedAt       time.Time      `json:"created_at"`
}

type KBSyncSource struct {
	ID               int64          `json:"id"`
	ChatbotID        int64          `json:"chatbot_id"`
	SourceType       string         `json:"source_type"`
	SourceIdentifier string         `json:"source_identifier"`
	AuthToken        sql.NullString `json:"auth_token,omitempty"`
	SyncFrequency    string         `json:"sync_frequency"`
	LastSyncAt       sql.NullTime   `json:"last_sync_at"`
	NextSyncAt       sql.NullTime   `json:"next_sync_at"`
	IsActive         bool           `json:"is_active"`
	SyncSettings     sql.NullString `json:"sync_settings"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type KBQualityTest struct {
	ID              int64           `json:"id"`
	ChatbotID       int64           `json:"chatbot_id"`
	TestQuery       string          `json:"test_query"`
	ExpectedContent sql.NullString  `json:"expected_content"`
	ActualResponse  sql.NullString  `json:"actual_response"`
	RelevanceScore  sql.NullFloat64 `json:"relevance_score"`
	Passed          sql.NullBool    `json:"passed"`
	TestedAt        time.Time       `json:"tested_at"`
}

type KBSearchIndex struct {
	ID        int64     `json:"id"`
	KBID      int64     `json:"kb_id"`
	Keyword   string    `json:"keyword"`
	Frequency int       `json:"frequency"`
	CreatedAt time.Time `json:"created_at"`
}

type Conversation struct {
	ID        int64        `json:"id"`
	ChatbotID int64        `json:"chatbot_id"`
	SessionID string       `json:"session_id"`
	VisitorID string       `json:"visitor_id"`
	StartedAt time.Time    `json:"started_at"`
	EndedAt   sql.NullTime `json:"ended_at"`
}

type Message struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
}

type User struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

type APIUsage struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	TokensUsed     int       `json:"tokens_used"`
	Cost           float64   `json:"cost"`
	Timestamp      time.Time `json:"timestamp"`
}

// LeadCaptureField represents a single form field definition.
type LeadCaptureField struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"`        // text | email | tel | textarea
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
}

type LeadCaptureConfig struct {
	ID        int64     `json:"id"`
	ChatbotID int64     `json:"chatbot_id"`
	Enabled   bool      `json:"enabled"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Fields    string    `json:"fields"` // JSON array of LeadCaptureField
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Lead struct {
	ID          int64     `json:"id"`
	ChatbotID   int64     `json:"chatbot_id"`
	SessionID   string    `json:"session_id"`
	FieldValues string    `json:"field_values"` // JSON object {"name":"...", "email":"..."}
	CreatedAt   time.Time `json:"created_at"`
}
