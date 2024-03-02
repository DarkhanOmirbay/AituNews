package postgresql

import (
	"darkhanomirbay.net/aitunews/pkg/models"
	"database/sql"
	"errors"
)

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) Insert(title, content, category string) (int, error) {
	stmt := `INSERT INTO articles (title, content, category) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, title, content, category).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ArticleModel) Get(id int) (*models.Article, error) {
	stmt := `SELECT id,title,content,category From articles Where id=$1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Article{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *ArticleModel) GetAll() ([]*models.Article, error) {
	stmt := `SELECT * from articles`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := []*models.Article{}

	for rows.Next() {
		a := &models.Article{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		articles = append(articles, a)

	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return articles, nil
}

func (m *ArticleModel) GetStudents() ([]*models.Article, error) {
	stmt := "SELECT * FROM articles Where category='For students'"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*models.Article{}

	for rows.Next() {
		a := &models.Article{}

		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return articles, nil
}
func (m *ArticleModel) GetTeachers() ([]*models.Article, error) {
	stmt := "SELECT * FROM articles Where category='For teachers'"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []*models.Article{}

	for rows.Next() {
		a := &models.Article{}

		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}
