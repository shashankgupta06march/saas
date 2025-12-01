# Advanced Knowledge Base API Documentation

## Overview
This document details all the new API endpoints for advanced knowledge base features.

**Base URL:** `http://localhost:8081/api`
**Authentication:** All endpoints require Bearer token in `Authorization` header

---

## 📁 Folder Management

### Create Folder
```http
POST /folders
Content-Type: application/json
Authorization: Bearer {token}

{
  "chatbot_id": 1,
  "name": "Product Documentation",
  "description": "All product docs",
  "parent_id": null,
  "color": "#3B82F6",
  "icon": "folder"
}
```

**Response 201:**
```json
{
  "id": 1,
  "chatbot_id": 1,
  "name": "Product Documentation",
  "description": "All product docs",
  "parent_id": null,
  "color": "#3B82F6",
  "icon": "folder",
  "created_at": "2025-10-15T12:00:00Z",
  "updated_at": "2025-10-15T12:00:00Z"
}
```

### Get Folder
```http
GET /folders/:id
Authorization: Bearer {token}
```

### Update Folder
```http
PUT /folders/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Updated Name",
  "description": "Updated description",
  "parent_id": 2,
  "color": "#EF4444",
  "icon": "book"
}
```

### Delete Folder
```http
DELETE /folders/:id
Authorization: Bearer {token}
```

**Note:** Cannot delete folder with children

### Get Folder Tree
```http
GET /folders/tree/:chatbot_id
Authorization: Bearer {token}
```

**Response 200:**
```json
[
  {
    "id": 1,
    "name": "Root Folder",
    "children": [
      {
        "id": 2,
        "name": "Child Folder",
        "children": []
      }
    ]
  }
]
```

### Get Root Folders
```http
GET /folders/roots/:chatbot_id
Authorization: Bearer {token}
```

### Move Folder
```http
POST /folders/:id/move
Content-Type: application/json
Authorization: Bearer {token}

{
  "new_parent_id": 3
}
```

**Use `"new_parent_id": null` to move to root level**

### Get Folder Path
```http
GET /folders/:id/path
Authorization: Bearer {token}
```

**Response 200:**
```json
[
  {"id": 1, "name": "Root"},
  {"id": 2, "name": "Parent"},
  {"id": 3, "name": "Current"}
]
```

---

## 🏷️ Tag Management

### Create Tag
```http
POST /tags
Content-Type: application/json
Authorization: Bearer {token}

{
  "chatbot_id": 1,
  "name": "important",
  "color": "#EF4444"
}
```

**Note:** Tag names are automatically normalized to lowercase

### Get Tag
```http
GET /tags/:id
Authorization: Bearer {token}
```

### Update Tag
```http
PUT /tags/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "urgent",
  "color": "#F59E0B"
}
```

### Delete Tag
```http
DELETE /tags/:id
Authorization: Bearer {token}
```

### Get All Tags for Chatbot
```http
GET /tags/chatbot/:chatbot_id
Authorization: Bearer {token}
```

### Assign Tags to Knowledge Entry
```http
POST /knowledge/:kb_id/tags
Content-Type: application/json
Authorization: Bearer {token}

{
  "tag_ids": [1, 2, 3]
}
```

**Note:** This replaces all existing tags

### Get Tags for Knowledge Entry
```http
GET /knowledge/:kb_id/tags
Authorization: Bearer {token}
```

### Add Single Tag to Knowledge Entry
```http
POST /knowledge/:kb_id/tags/:tag_id
Authorization: Bearer {token}
```

### Remove Tag from Knowledge Entry
```http
DELETE /knowledge/:kb_id/tags/:tag_id
Authorization: Bearer {token}
```

---

## 📜 Version Control

### Get Version History
```http
GET /knowledge/:kb_id/versions
Authorization: Bearer {token}
```

**Response 200:**
```json
[
  {
    "id": 10,
    "kb_id": 1,
    "version": 3,
    "title": "Document Title",
    "content": "Latest content...",
    "content_type": "text",
    "source_url": "",
    "changed_by": 5,
    "change_summary": "Updated formatting",
    "created_at": "2025-10-15T14:00:00Z"
  },
  {
    "id": 8,
    "kb_id": 1,
    "version": 2,
    ...
  }
]
```

### Get Specific Version
```http
GET /knowledge/:kb_id/versions/:version
Authorization: Bearer {token}
```

### Restore to Previous Version
```http
POST /knowledge/:kb_id/versions/:version/restore
Authorization: Bearer {token}
```

**Behavior:**
- Saves current state as new version before restoring
- Creates new version with restored content
- Increments version number

### Compare Two Versions
```http
GET /knowledge/:kb_id/versions/compare?version1=1&version2=3
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "version1": 1,
  "version2": 3,
  "changes": ["title", "content"],
  "has_changes": true
}
```

---

## 🧩 Smart Chunking

### Get Chunks for Knowledge Entry
```http
GET /knowledge/:kb_id/chunks
Authorization: Bearer {token}
```

**Response 200:**
```json
[
  {
    "id": 1,
    "kb_id": 1,
    "chunk_index": 0,
    "content": "First chunk content...",
    "embedding_vector": "[0.1, 0.2, ...]",
    "token_count": 150,
    "metadata": "{\"page\": 1}",
    "created_at": "2025-10-15T12:00:00Z"
  },
  ...
]
```

### Rechunk Knowledge Entry
```http
POST /knowledge/:kb_id/rechunk
Content-Type: application/json
Authorization: Bearer {token}

{
  "chunk_size": 1500,
  "chunk_overlap": 300,
  "chunking_method": "semantic"
}
```

**Chunking Methods:**
- `"fixed"` - Fixed-size chunks with overlap
- `"sentence"` - Sentence-boundary aware
- `"semantic"` - Paragraph-aware splitting

**All parameters are optional. Existing KB settings used if not provided.**

**Response 200:**
```json
{
  "message": "Knowledge base rechunked successfully",
  "chunk_count": 8,
  "chunks": [...]
}
```

### Preview Chunking
```http
POST /knowledge/preview-chunks
Content-Type: application/json
Authorization: Bearer {token}

{
  "content": "Long text to be chunked...",
  "chunk_size": 1000,
  "chunk_overlap": 200,
  "chunking_method": "semantic"
}
```

**Response 200:**
```json
{
  "chunk_count": 5,
  "chunks": ["chunk 1...", "chunk 2...", ...],
  "chunk_size": 1000,
  "chunk_overlap": 200,
  "chunking_method": "semantic"
}
```

---

## 🧪 Quality Testing

### Run Single Test
```http
POST /quality/test
Content-Type: application/json
Authorization: Bearer {token}

{
  "chatbot_id": 1,
  "test_query": "How do I reset my password?",
  "expected_content": "password reset"
}
```

**Response 200:**
```json
{
  "id": 1,
  "chatbot_id": 1,
  "test_query": "How do I reset my password?",
  "expected_content": "password reset",
  "actual_response": "To reset your password, go to...",
  "relevance_score": 0.85,
  "passed": true,
  "tested_at": "2025-10-15T12:00:00Z"
}
```

**Pass/Fail Threshold:** 0.6 (60% relevance)

### Run Batch Tests
```http
POST /quality/batch-test
Content-Type: application/json
Authorization: Bearer {token}

{
  "chatbot_id": 1,
  "tests": [
    {
      "query": "What is your return policy?",
      "expected": "30 days"
    },
    {
      "query": "Do you ship internationally?",
      "expected": "worldwide shipping"
    }
  ]
}
```

**Response 200:**
```json
{
  "total_tests": 2,
  "results": [...]
}
```

### Get Test History
```http
GET /quality/tests/:chatbot_id?limit=50
Authorization: Bearer {token}
```

### Get Test Statistics
```http
GET /quality/stats/:chatbot_id
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "total_tests": 100,
  "passed_tests": 75,
  "failed_tests": 25,
  "average_score": 0.72,
  "pass_rate": 75.0
}
```

### Get Failed Tests
```http
GET /quality/failed/:chatbot_id?limit=20
Authorization: Bearer {token}
```

### Identify Knowledge Gaps
```http
GET /quality/gaps/:chatbot_id
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "gaps": ["password", "shipping", "refund", "payment", "account"]
}
```

**Returns keywords from failed test queries**

### Get Quality Score
```http
GET /quality/score/:chatbot_id
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "chatbot_id": 1,
  "quality_score": 0.78,
  "rating": "Good"
}
```

**Rating Scale:**
- `0.9+` - Excellent
- `0.75-0.89` - Good
- `0.6-0.74` - Fair
- `0.4-0.59` - Poor
- `<0.4` - Very Poor

### Get Improvement Suggestions
```http
GET /quality/suggestions/:chatbot_id
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "suggestions": [
    "Pass rate is low. Review failed tests and add missing content to knowledge base.",
    "Common topics in failed tests: password, shipping, refund"
  ]
}
```

---

## 🔄 Sync Source Management

### Create Sync Source
```http
POST /sync-sources
Content-Type: application/json
Authorization: Bearer {token}

{
  "chatbot_id": 1,
  "source_type": "google_drive",
  "source_identifier": "folder_id_or_url",
  "auth_token": "oauth_token",
  "sync_frequency": "daily",
  "sync_settings": "{\"filter\": \".pdf\"}"
}
```

**Source Types:**
- `google_drive`
- `dropbox`
- `notion`
- `wordpress`
- `confluence`
- `webhook`

**Sync Frequencies:**
- `realtime` (checks every 5 minutes)
- `hourly`
- `daily`
- `weekly`

### Get Sync Source
```http
GET /sync-sources/:id
Authorization: Bearer {token}
```

### Update Sync Source
```http
PUT /sync-sources/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "source_type": "google_drive",
  "source_identifier": "new_folder_id",
  "auth_token": "new_token",
  "sync_frequency": "hourly",
  "sync_settings": "{}"
}
```

### Delete Sync Source
```http
DELETE /sync-sources/:id
Authorization: Bearer {token}
```

### Get All Sync Sources for Chatbot
```http
GET /sync-sources/chatbot/:chatbot_id
Authorization: Bearer {token}
```

### Toggle Sync Source (Enable/Disable)
```http
POST /sync-sources/:id/toggle
Content-Type: application/json
Authorization: Bearer {token}

{
  "active": true
}
```

### Trigger Manual Sync
```http
POST /sync-sources/:id/trigger
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "message": "Sync triggered successfully"
}
```

### Get Sync Status
```http
GET /sync-sources/:id/status
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "id": 1,
  "source_type": "google_drive",
  "is_active": true,
  "sync_frequency": "daily",
  "last_synced_at": "2025-10-14T00:00:00Z",
  "next_sync_at": "2025-10-15T00:00:00Z",
  "sync_overdue": false
}
```

---

## 📝 Enhanced Knowledge Base Endpoints

### Upload File (Enhanced)
```http
POST /knowledge/upload
Content-Type: multipart/form-data
Authorization: Bearer {token}

file: [PDF, DOCX, TXT, XLSX, XLS, CSV file]
chatbot_id: 1
title: Document Title
```

**Now Supports:**
- PDF files
- Word documents (.docx)
- Text files (.txt)
- Excel files (.xlsx, .xls) ✨ **NEW**
- CSV files (.csv) ✨ **NEW**

**Response 200:**
```json
{
  "id": 1,
  "title": "Document Title",
  "content_type": "xlsx",
  "message": "Knowledge added successfully"
}
```

---

## 📊 Example Workflows

### Workflow 1: Organize Knowledge with Folders & Tags

```bash
# 1. Create a folder
curl -X POST http://localhost:8081/api/folders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chatbot_id": 1,
    "name": "Product Docs",
    "color": "#3B82F6"
  }'

# 2. Create tags
curl -X POST http://localhost:8081/api/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chatbot_id": 1,
    "name": "important",
    "color": "#EF4444"
  }'

# 3. Assign tags to knowledge entry
curl -X POST http://localhost:8081/api/knowledge/1/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tag_ids": [1, 2]
  }'

# 4. Get folder tree
curl http://localhost:8081/api/folders/tree/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 2: Quality Testing & Improvement

```bash
# 1. Run quality tests
curl -X POST http://localhost:8081/api/quality/batch-test \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chatbot_id": 1,
    "tests": [
      {"query": "How to reset password?", "expected": "password reset"},
      {"query": "Shipping policy?", "expected": "free shipping"}
    ]
  }'

# 2. Get statistics
curl http://localhost:8081/api/quality/stats/1 \
  -H "Authorization: Bearer $TOKEN"

# 3. Identify gaps
curl http://localhost:8081/api/quality/gaps/1 \
  -H "Authorization: Bearer $TOKEN"

# 4. Get suggestions
curl http://localhost:8081/api/quality/suggestions/1 \
  -H "Authorization: Bearer $TOKEN"

# 5. Get quality score
curl http://localhost:8081/api/quality/score/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 3: Version Control

```bash
# 1. Get version history
curl http://localhost:8081/api/knowledge/1/versions \
  -H "Authorization: Bearer $TOKEN"

# 2. Compare versions
curl http://localhost:8081/api/knowledge/1/versions/compare?version1=1&version2=3 \
  -H "Authorization: Bearer $TOKEN"

# 3. Restore to previous version
curl -X POST http://localhost:8081/api/knowledge/1/versions/2/restore \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 4: Smart Chunking

```bash
# 1. Preview chunking
curl -X POST http://localhost:8081/api/knowledge/preview-chunks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Long text content...",
    "chunk_size": 1000,
    "chunk_overlap": 200,
    "chunking_method": "semantic"
  }'

# 2. Rechunk existing entry
curl -X POST http://localhost:8081/api/knowledge/1/rechunk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chunk_size": 1500,
    "chunking_method": "semantic"
  }'

# 3. View chunks
curl http://localhost:8081/api/knowledge/1/chunks \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔒 Authentication

All endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Get token from login:
```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'
```

---

## 📈 Response Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## 💡 Tips

1. **Folders:** Use colors and icons for visual organization
2. **Tags:** Keep tag names short and consistent
3. **Chunking:** Start with `semantic` method for best results
4. **Quality Tests:** Run tests regularly (weekly recommended)
5. **Version Control:** Add change summaries for clarity
6. **Sync Sources:** Start with `daily` frequency, adjust as needed

---

**API Version:** 1.0
**Last Updated:** October 15, 2025
**Backend Port:** 8081



