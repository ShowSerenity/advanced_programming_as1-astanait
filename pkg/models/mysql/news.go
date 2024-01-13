package mysql

import (
	"database/sql"
	"errors"
	"showserenity.net/astanait/pkg/models"
)

type NewsModel struct {
	DB *sql.DB
}

func (m *NewsModel) Insert(title, content, expires, newsType string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires, type)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`
	result, err := m.DB.Exec(stmt, title, content, expires, newsType)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *NewsModel) GetById(id int) (*models.News, error) {
	stmt := `SELECT id, title, content, created, expires, type FROM news
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.News{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *NewsModel) GetByType(newsType string) ([]*models.News, error) {
	stmt := `SELECT id, title, content, created, expires, type FROM news
	WHERE expires > UTC_TIMESTAMP() AND type = ?`

	rows, err := m.DB.Query(stmt, newsType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []*models.News{}

	for rows.Next() {
		s := &models.News{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Type)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, s)
	}

	return newsList, nil
}

func (m *NewsModel) Latest() ([]*models.News, error) {
	stmt := `SELECT id, title, content, created, expires, type FROM news
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := []*models.News{}

	for rows.Next() {
		s := &models.News{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Type)
		if err != nil {
			return nil, err
		}

		news = append(news, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}
