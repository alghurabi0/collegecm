package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateUser(v *validator.Validator, user *User) {
	// TODO - handle strings length with varchar
	v.Check(user.Username != "", "اسم الحساب", "يجب تزويد المعلومات")
	v.Check(user.Password != "", "الرمز", "يجب تزويد المعلومات")
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (username, password) 
        VALUES ($1, $2)
        RETURNING id, created_at`
	args := []interface{}{
		user.Username,
		user.Password,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `SELECT * FROM %s users`
	var args []interface{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserModel) Get(id int64) (*User, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT * FROM users WHERE id = $1;`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) Update(user *User) error {
	if user.ID < 0 {
		return ErrRecordNotFound
	}
	query := `
	UPDATE users
	SET username = $1, password = $2
	WHERE student_id = $3`
	args := []interface{}{
		&user.Username,
		&user.Password,
		&user.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(ctx, query, args...)
	return err
}

func (u UserModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (u UserModel) GetByUsername(username string) (*User, error) {
	query := `
	SELECT * FROM users WHERE username = $1;`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
