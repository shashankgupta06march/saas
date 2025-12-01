package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBQualityRepository struct {
	db *sql.DB
}

func NewKBQualityRepository(db *sql.DB) *KBQualityRepository {
	return &KBQualityRepository{db: db}
}

func (r *KBQualityRepository) Create(test *models.KBQualityTest) error {
	query := `INSERT INTO kb_quality_tests 
	          (chatbot_id, test_query, expected_content, actual_response, relevance_score, passed) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, test.ChatbotID, test.TestQuery, test.ExpectedContent,
		test.ActualResponse, test.RelevanceScore, test.Passed)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	test.ID = id
	return nil
}

func (r *KBQualityRepository) GetByID(id int64) (*models.KBQualityTest, error) {
	query := `SELECT id, chatbot_id, test_query, expected_content, actual_response, 
	          relevance_score, passed, tested_at 
	          FROM kb_quality_tests WHERE id = ?`

	test := &models.KBQualityTest{}
	err := r.db.QueryRow(query, id).Scan(
		&test.ID, &test.ChatbotID, &test.TestQuery, &test.ExpectedContent,
		&test.ActualResponse, &test.RelevanceScore, &test.Passed, &test.TestedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (r *KBQualityRepository) GetByChatbot(chatbotID int64, limit int) ([]models.KBQualityTest, error) {
	query := `SELECT id, chatbot_id, test_query, expected_content, actual_response, 
	          relevance_score, passed, tested_at 
	          FROM kb_quality_tests WHERE chatbot_id = ? 
	          ORDER BY tested_at DESC LIMIT ?`

	rows, err := r.db.Query(query, chatbotID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.KBQualityTest
	for rows.Next() {
		var test models.KBQualityTest
		err := rows.Scan(
			&test.ID, &test.ChatbotID, &test.TestQuery, &test.ExpectedContent,
			&test.ActualResponse, &test.RelevanceScore, &test.Passed, &test.TestedAt,
		)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func (r *KBQualityRepository) GetStats(chatbotID int64) (map[string]interface{}, error) {
	query := `SELECT 
	          COUNT(*) as total_tests,
	          SUM(CASE WHEN passed = true THEN 1 ELSE 0 END) as passed_tests,
	          AVG(relevance_score) as avg_score
	          FROM kb_quality_tests 
	          WHERE chatbot_id = ?`

	stats := make(map[string]interface{})
	var totalTests, passedTests int
	var avgScore sql.NullFloat64

	err := r.db.QueryRow(query, chatbotID).Scan(&totalTests, &passedTests, &avgScore)
	if err != nil {
		return nil, err
	}

	stats["total_tests"] = totalTests
	stats["passed_tests"] = passedTests
	stats["failed_tests"] = totalTests - passedTests

	if avgScore.Valid {
		stats["average_score"] = avgScore.Float64
	} else {
		stats["average_score"] = 0.0
	}

	if totalTests > 0 {
		stats["pass_rate"] = float64(passedTests) / float64(totalTests) * 100
	} else {
		stats["pass_rate"] = 0.0
	}

	return stats, nil
}

func (r *KBQualityRepository) Delete(id int64) error {
	query := `DELETE FROM kb_quality_tests WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *KBQualityRepository) DeleteByChatbot(chatbotID int64) error {
	query := `DELETE FROM kb_quality_tests WHERE chatbot_id = ?`
	_, err := r.db.Exec(query, chatbotID)
	return err
}

func (r *KBQualityRepository) GetFailedTests(chatbotID int64, limit int) ([]models.KBQualityTest, error) {
	query := `SELECT id, chatbot_id, test_query, expected_content, actual_response, 
	          relevance_score, passed, tested_at 
	          FROM kb_quality_tests 
	          WHERE chatbot_id = ? AND passed = false 
	          ORDER BY tested_at DESC LIMIT ?`

	rows, err := r.db.Query(query, chatbotID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.KBQualityTest
	for rows.Next() {
		var test models.KBQualityTest
		err := rows.Scan(
			&test.ID, &test.ChatbotID, &test.TestQuery, &test.ExpectedContent,
			&test.ActualResponse, &test.RelevanceScore, &test.Passed, &test.TestedAt,
		)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}


