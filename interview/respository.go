package interview

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) CreateInterview(interview *Interview) (int, error) {
	query := `
    INSERT INTO interviews (user_id, length, number_questions, difficulty, status, score, language, first_question)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id
    `

	var id int
	err := repo.DB.QueryRow(query,
		interview.UserId,
		interview.Length,
		interview.NumberQuestions,
		interview.Difficulty,
		interview.Status,
		interview.Score,
		interview.Language,
		interview.FirstQuestion).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetInterview(interviewID int) (*Interview, error) {
	query := `
	SELECT (user_id, length, number_questions, difficulty, status, score, language, first_question)
	FROM interviews
	WHERE interview_id = $1
	`

	interview := &Interview{}
	err := repo.DB.QueryRow(query,
		interviewID).Scan(
		&interview.UserId,
		&interview.Length,
		&interview.NumberQuestions,
		&interview.Difficulty,
		&interview.Status,
		&interview.Score,
		&interview.Language,
		&interview.FirstQuestion)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no interview found with id %d", interviewID)
		}
		return nil, err
	}

	return interview, nil
}
