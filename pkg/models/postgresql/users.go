package postgresql

import (
	"darkhanomirbay.net/aitunews/pkg/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, role string) error {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (full_name,email,hashed_password,role,created)
values($1,$2,$3,$4,NOW())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedpassword), role)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "users_email_key") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 AND active = true"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

func (m *UserModel) GetAllUsers() ([]*models.User, error) {
	stmt := `SELECT id,full_name,email,role from users`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		u := &models.User{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		users = append(users, u)

	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return users, nil
}
func (m *UserModel) IsAdmin(id int) (string, error) {
	stmt := `SELECT role FROM users WHERE id=$1`
	row := m.DB.QueryRow(stmt, id)

	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNoRecord
		} else {
			return "", err
		}
	}
	return role, nil
}
func (m *UserModel) IsTeacher(id int) (string, error) {
	stmt := `SELECT role From users where id=$1`

	row := m.DB.QueryRow(stmt, id)

	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNoRecord
		} else {
			return "", err
		}
	}
	return role, err

}
