package db

import (
	"context"
)

const createLemburQuery = `INSERT INTO lembur (id, no_payroll, nama, jabatan, gaji_pokok, basis) VALUES (?,?, ?, ?, ?, ?)`

type CreateLemburParams struct {
	ID        string  `json:"id"`
	NoPayroll string  `json:"no_payroll"`
	Nama      string  `json:"nama"`
	Jabatan   string  `json:"jabatan"`
	GajiPokok float32 `json:"gaji_pokok"`
	Basis     float32 `json:"basis"`
}

func (q *Queries) CreateLembur(ctx context.Context, arg CreateLemburParams) error {
	_, err := q.db.ExecContext(ctx, createLemburQuery,
		arg.ID,
		arg.NoPayroll,
		arg.Nama,
		arg.Jabatan,
		arg.GajiPokok,
		arg.Basis,
	)
	return err
}

const createTanggalLemburQuery = `INSERT INTO tanggal_lembur (id, tanggal, total_jam, id_lembur) VALUES (?, ?, ?, ?)`

type CreateTanggalLemburParams struct {
	ID       string  `json:"id"`
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
	IDLembur string  `json:"id_lembur"`
}

func (q *Queries) CreateTanggalLembur(ctx context.Context, arg CreateTanggalLemburParams) error {
	_, err := q.db.ExecContext(ctx, createTanggalLemburQuery, arg.ID, arg.Tanggal, arg.TotalJam, arg.IDLembur)
	return err
}

const getLemburByIDQuery = `SELECT id, no_payroll, nama, jabatan, gaji_pokok, basis FROM lembur WHERE id = ? LIMIT 1`

func (q *Queries) GetLemburByID(ctx context.Context, id string) (Lembur, error) {
	row := q.db.QueryRowContext(ctx, getLemburByIDQuery, id)
	var i Lembur
	err := row.Scan(
		&i.ID,
		&i.NoPayroll,
		&i.Nama,
		&i.Jabatan,
		&i.GajiPokok,
		&i.Basis,
	)
	return i, err
}

const GetLemburByNoPayrollQuery = `SELECT id, no_payroll, nama, jabatan, gaji_pokok, basis FROM lembur WHERE no_payroll = ? LIMIT 1`

func (q *Queries) GetLemburByNoPayroll(ctx context.Context, noPayroll string) (Lembur, error) {
	row := q.db.QueryRowContext(ctx, GetLemburByNoPayrollQuery, noPayroll)
	var i Lembur
	err := row.Scan(
		&i.ID,
		&i.NoPayroll,
		&i.Nama,
		&i.Jabatan,
		&i.GajiPokok,
		&i.Basis,
	)
	return i, err
}

const getTanggalLemburByIdLemburQuery = `SELECT id, tanggal, total_jam FROM tanggal_lembur WHERE id_lembur = ?`

type GetTanggalLemburByIdLemburRow struct {
	ID       string  `json:"id"`
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
}

func (q *Queries) GetTanggalLemburByIdLembur(ctx context.Context, idLembur string) ([]GetTanggalLemburByIdLemburRow, error) {
	rows, err := q.db.QueryContext(ctx, getTanggalLemburByIdLemburQuery, idLembur)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetTanggalLemburByIdLemburRow
	for rows.Next() {
		var i GetTanggalLemburByIdLemburRow
		if err := rows.Scan(&i.ID, &i.Tanggal, &i.TotalJam); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLemburQuery = `SELECT id, no_payroll, nama, jabatan, gaji_pokok, basis FROM lembur`

func (q *Queries) ListLembur(ctx context.Context) ([]Lembur, error) {
	rows, err := q.db.QueryContext(ctx, listLemburQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Lembur
	for rows.Next() {
		var i Lembur
		if err := rows.Scan(
			&i.ID,
			&i.NoPayroll,
			&i.Nama,
			&i.Jabatan,
			&i.GajiPokok,
			&i.Basis,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLemburWithTanggalLemburQuery = `
SELECT
    lembur.id,
    lembur.no_payroll,
    lembur.nama,
    lembur.jabatan,
    lembur.gaji_pokok,
    lembur.basis,
    tanggal_lembur.id AS tanggal_lembur_id,
    tanggal_lembur.tanggal,
    tanggal_lembur.total_jam
FROM
    lembur
INNER JOIN
    tanggal_lembur ON lembur.id = tanggal_lembur.id_lembur
`

type ListLemburWithTanggalLemburRow struct {
	ID              string  `json:"id"`
	NoPayroll       string  `json:"no_payroll"`
	Nama            string  `json:"nama"`
	Jabatan         string  `json:"jabatan"`
	GajiPokok       float32 `json:"gaji_pokok"`
	Basis           float32 `json:"basis"`
	TanggalLemburID string  `json:"tanggal_lembur_id"`
	Tanggal         string  `json:"tanggal"`
	TotalJam        float32 `json:"total_jam"`
}

func (q *Queries) ListLemburWithTanggalLembur(ctx context.Context) ([]ListLemburWithTanggalLemburRow, error) {
	rows, err := q.db.QueryContext(ctx, listLemburWithTanggalLemburQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLemburWithTanggalLemburRow{}
	for rows.Next() {
		var i ListLemburWithTanggalLemburRow
		if err := rows.Scan(
			&i.ID,
			&i.NoPayroll,
			&i.Nama,
			&i.Jabatan,
			&i.GajiPokok,
			&i.Basis,
			&i.TanggalLemburID,
			&i.Tanggal,
			&i.TotalJam,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTanggalLemburQuery = `SELECT id, id_lembur, tanggal, total_jam FROM tanggal_lembur`

type ListTanggalLemburRow struct {
	ID       string  `json:"id"`
	IDLembur string  `json:"id_lembur"`
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
}

func (q *Queries) ListTanggalLembur(ctx context.Context) ([]ListTanggalLemburRow, error) {
	rows, err := q.db.QueryContext(ctx, listTanggalLemburQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTanggalLemburRow

	for rows.Next() {
		var i ListTanggalLemburRow
		if err := rows.Scan(
			&i.ID,
			&i.IDLembur,
			&i.Tanggal,
			&i.TotalJam,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) DeleteLembur(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, `DELETE FROM lembur WHERE id = ?`, id)
	return err
}
