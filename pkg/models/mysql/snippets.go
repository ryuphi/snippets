package mysql

import (
	"database/sql"
	"errors"
	"learn-web/snippets/pkg/models"
)

type SnippetModel struct {
	Db *sql.DB
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// write the sql statement we want to execute.
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec method on the embedded connection pool to execute the statement.
	result, err := m.Db.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId method on the result object to get the Id of our newly inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// the Id has the type int64, so we convert it to an int type
	return int(id), nil
}

// Get return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}

	err := m.Db.QueryRow(`SELECT 
       		id,
			title,
       		content,
       		created,
       		expires 
		FROM snippets 
		WHERE expires > UTC_TIMESTAMP() AND id = ?`,
		id).Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// all ok, return the Snippet object
	return s, nil
}

// Latest return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.Db.Query(stmt)
	if err != nil {
		return nil, err
	}

	// closing the resulset.
	defer rows.Close()

	// initialize an empty slice to hold the models.Snippet objects
	var snippets []*models.Snippet

	// iterate through the rows in the resultset.
	// when finish the iteration over all the rows, the resultset automatically closes itself and free-up the
	// database connection.
	for rows.Next() {
		s := &models.Snippet{}

		// rows.Scan copy the values from each field in the row to the new Snippet object that we created.
		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// append it to the slice...
		snippets = append(snippets, s)
	}

	// if everything is ok, return the snippets slice.
	return snippets, nil
}
