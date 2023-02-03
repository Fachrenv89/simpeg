package db

import (
	"context"
)

type Querier interface {
	CreateLembur(ctx context.Context, arg CreateLemburParams) error
	CreateTanggalLembur(ctx context.Context, arg CreateTanggalLemburParams) error
	GetLemburByID(ctx context.Context, id string) (Lembur, error)
	GetLemburByNoPayroll(ctx context.Context, noPayroll string) (Lembur, error)
	GetTanggalLemburByIdLembur(ctx context.Context, idLembur string) ([]GetTanggalLemburByIdLemburRow, error)
	ListLembur(ctx context.Context) ([]Lembur, error)
	ListLemburWithTanggalLembur(ctx context.Context) ([]ListLemburWithTanggalLemburRow, error)
	ListTanggalLembur(ctx context.Context) ([]ListTanggalLemburRow, error)
	DeleteLembur(ctx context.Context, id string) error

	ListNonStaff(ctx context.Context) ([]NonStaff, error)
	GetNonStaffByNoPayroll(ctx context.Context, noPayroll string) (NonStaff, error)
	GetNonStaffByID(ctx context.Context, id string) (NonStaff, error)
	CreateNonStaff(ctx context.Context, arg CreateNonStaffParams) error
	CountNonStaff(ctx context.Context) int
	DeleteNonStaff(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) error
	UpdateNonStaff(ctx context.Context, arg UpdateNonStaffParams) error
	UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
}

var _ Querier = (*Queries)(nil)
