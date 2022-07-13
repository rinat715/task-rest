package users

import (
	"database/sql"
	e "go_rest/internal/errors"
	m "go_rest/internal/models"
)

var baseQuery string = `SELECT userid, email, pass, is_admin
						FROM users`

type UserRepositoryInterface interface {
	Create(user *m.User) error
	Get(userId int) (m.User, error)
	GetbyEmail(email string) (m.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (s *UserRepository) Close() error {
	return s.db.Close()
}

func (s *UserRepository) Create(user *m.User) error {
	var err error
	query := "INSERT INTO users(email, pass, is_admin) VALUES(?, ?, ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(user.Email, user.Pass, user.IsAdmin)
	if err != nil {
		return err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = int(userId)
	return nil

}

func (s *UserRepository) Get(userId int) (m.User, error) {
	var user m.User
	query := baseQuery + `\n` + `WHERE userid = ?`

	err := s.db.QueryRow(query, userId).Scan(&user.Id, &user.Email, &user.Pass, &user.IsAdmin)
	switch {
	default:
		return user, nil
	case err == sql.ErrNoRows:
		return user, &e.UserNotFound{UserId: userId}
	case err != nil:
		return user, err
	}
}

func (s *UserRepository) GetbyEmail(email string) (m.User, error) {
	var user m.User
	query := baseQuery + `\n` + `WHERE email = ?`

	err := s.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Pass, &user.IsAdmin)
	switch {
	default:
		return user, nil
	case err == sql.ErrNoRows:
		return user, &e.UserNotFound{Email: email}
	case err != nil:
		return user, err
	}
}
