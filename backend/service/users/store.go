package users

import (
	"database/sql"
	"fmt"

	"github.com/alexhool2/TimeCard/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users(firstName, lastName, userName, password, email, role) VALUES (?,?,?,?,?,?)", user.FirstName, user.LastName, user.UserName, user.Password, user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAllUsers() ([]types.User, error) {
	rows, err := s.db.Query("SELECT id, firstName, lastName, userName, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []types.User
	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.UserName, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) GetDynamicUsers(search string) ([]types.User, error) {
	search = "%" + search + "%"
	rows, err := s.db.Query("SELECT id, firstName, lastName, userName, email from users WHERE firstName LIKE ? OR userName LIKE ? OR email LIKE ?", search, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []types.User
	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.UserName, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)

	}
	return users, nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) GetUserByUserName(userName string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE username = ?", userName)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users set password = ? WHERE id = ?`
	_, err := s.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password for user %d: %v", userID, err)
	}
	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Password, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
