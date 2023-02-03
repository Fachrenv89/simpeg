package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name, role, email, password, photo) VALUES (?, ?, ?, ?, ?, ?);`

type CreateUserParams struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Role,
		arg.Email,
		arg.Password,
		arg.Photo,
	)

	return err
}

type UpdateUserPasswordParams struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users SET password = ? WHERE id = ?;`

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword,
		arg.Password,
		arg.ID,
	)
	return err
}

type UpdateUserProfileParams struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

const updateUser = `UPDATE users SET name = ?,  role = ?, email = ?, photo = ? WHERE id = ?;`

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.Role,
		arg.Email,
		arg.Photo,
		arg.ID,
	)

	return err
}

const getUser = `-- name: GetUser :one
SELECT id, name, role, email, password, photo FROM users WHERE id = ? LIMIT 1`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Role,
		&i.Email,
		&i.Password,
		&i.Photo,
	)

	return i, err
}

const getUserByEmail = `-- name: GetUser :one
SELECT id, name, role, email, password, photo FROM users WHERE email = ? LIMIT 1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Role,
		&i.Email,
		&i.Password,
		&i.Photo,
	)
	return i, err
}
