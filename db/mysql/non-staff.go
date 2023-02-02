package db

import (
	"context"
)

const createNonStaff = `-- name: CreateNonStaff :exec
INSERT INTO non_staff (
    id,
    no_payroll,
    nama,
    foto,
    status_pph_21,
    alamat,
    nik,
    npwp,
    gender,
    jabatan,
    golongan,
    expat,
    ptkp,
    status_kerja,
    awal_kerja,
    akhir_kerja,
    gaji_pokok,
    tunjangan_kemahalan,
    tunjangan_perumahan,
    tunjangan_jabatan,
    tunjangan_lain_pph21
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

type CreateNonStaffParams struct {
	ID                 string   `json:"id"`
	NoPayroll          string   `json:"no_payroll"`
	Nama               string   `json:"nama"`
	Foto               string   `json:"foto"`
	StatusPph21        string   `json:"status_pph_21"`
	Alamat             string   `json:"alamat"`
	Nik                string   `json:"nik"`
	Npwp               string   `json:"npwp"`
	Gender             string   `json:"gender"`
	Jabatan            string   `json:"jabatan"`
	Golongan           string   `json:"golongan"`
	Expat              string   `json:"expat"`
	Ptkp               string   `json:"ptkp"`
	StatusKerja        string   `json:"status_kerja"`
	AwalKerja          int32    `json:"awal_kerja"`
	AkhirKerja         int32    `json:"akhir_kerja"`
	GajiPokok          float32  `json:"gaji_pokok"`
	TunjanganKemahalan *float32 `json:"tunjangan_kemahalan"`
	TunjanganPerumahan *float32 `json:"tunjangan_perumahan"`
	TunjanganJabatan   *float32 `json:"tunjangan_jabatan"`
	TunjanganLainPph21 *float32 `json:"tunjangan_lain_pph21"`
}

type UpdateNonStaffParams struct {
	ID                 string   `json:"id"`
	NoPayroll          string   `json:"no_payroll"`
	Nama               string   `json:"nama"`
	Foto               string   `json:"foto"`
	StatusPph21        string   `json:"status_pph_21"`
	Alamat             string   `json:"alamat"`
	Nik                string   `json:"nik"`
	Npwp               string   `json:"npwp"`
	Gender             string   `json:"gender"`
	Jabatan            string   `json:"jabatan"`
	Golongan           string   `json:"golongan"`
	Expat              string   `json:"expat"`
	Ptkp               string   `json:"ptkp"`
	StatusKerja        string   `json:"status_kerja"`
	AwalKerja          int32    `json:"awal_kerja"`
	AkhirKerja         int32    `json:"akhir_kerja"`
	GajiPokok          float32  `json:"gaji_pokok"`
	TunjanganKemahalan *float32 `json:"tunjangan_kemahalan"`
	TunjanganPerumahan *float32 `json:"tunjangan_perumahan"`
	TunjanganJabatan   *float32 `json:"tunjangan_jabatan"`
	TunjanganLainPph21 *float32 `json:"tunjangan_lain_pph21"`
}

const updateNonStaff = `-- name: UpdateNonStaff :exec
UPDATE non_staff
SET 
	no_payroll = ?,
	nama = ?,
 	foto = ?,
    status_pph_21 = ?,
    alamat = ?,
    nik = ?,
    npwp = ?,
    gender = ?,
    jabatan = ?,
    golongan = ?,
    expat = ?,
    ptkp = ?,
    status_kerja = ?,
    awal_kerja = ?,
    akhir_kerja = ?,
    gaji_pokok = ?,
    tunjangan_kemahalan = ?,
    tunjangan_perumahan = ?,
    tunjangan_jabatan = ?,
    tunjangan_lain_pph21 = ?
WHERE id = ?
`

func (q *Queries) UpdateNonStaff(ctx context.Context, arg UpdateNonStaffParams) error {
	_, err := q.db.ExecContext(ctx, updateNonStaff,
		arg.NoPayroll,
		arg.Nama,
		arg.Foto,
		arg.StatusPph21,
		arg.Alamat,
		arg.Nik,
		arg.Npwp,
		arg.Gender,
		arg.Jabatan,
		arg.Golongan,
		arg.Expat,
		arg.Ptkp,
		arg.StatusKerja,
		arg.AwalKerja,
		arg.AkhirKerja,
		arg.GajiPokok,
		arg.TunjanganKemahalan,
		arg.TunjanganPerumahan,
		arg.TunjanganJabatan,
		arg.TunjanganLainPph21,
		arg.ID,
	)
	return err
}

func (q *Queries) CreateNonStaff(ctx context.Context, arg CreateNonStaffParams) error {
	_, err := q.db.ExecContext(ctx, createNonStaff,
		arg.ID,
		arg.NoPayroll,
		arg.Nama,
		arg.Foto,
		arg.StatusPph21,
		arg.Alamat,
		arg.Nik,
		arg.Npwp,
		arg.Gender,
		arg.Jabatan,
		arg.Golongan,
		arg.Expat,
		arg.Ptkp,
		arg.StatusKerja,
		arg.AwalKerja,
		arg.AkhirKerja,
		arg.GajiPokok,
		arg.TunjanganKemahalan,
		arg.TunjanganPerumahan,
		arg.TunjanganJabatan,
		arg.TunjanganLainPph21,
	)
	return err
}

const getNonStaffByID = `-- name: GetNonStaffByID :one
SELECT id, no_payroll, nama, foto, status_pph_21, alamat, nik, npwp, gender, jabatan, golongan, expat, ptkp, status_kerja, awal_kerja, akhir_kerja, gaji_pokok, tunjangan_kemahalan, tunjangan_perumahan, tunjangan_jabatan, tunjangan_lain_pph21 FROM non_staff WHERE id = ? LIMIT 1
`

func (q *Queries) GetNonStaffByID(ctx context.Context, id string) (NonStaff, error) {
	row := q.db.QueryRowContext(ctx, getNonStaffByID, id)
	var i NonStaff
	err := row.Scan(
		&i.ID,
		&i.NoPayroll,
		&i.Nama,
		&i.Foto,
		&i.StatusPph21,
		&i.Alamat,
		&i.Nik,
		&i.Npwp,
		&i.Gender,
		&i.Jabatan,
		&i.Golongan,
		&i.Expat,
		&i.Ptkp,
		&i.StatusKerja,
		&i.AwalKerja,
		&i.AkhirKerja,
		&i.GajiPokok,
		&i.TunjanganKemahalan,
		&i.TunjanganPerumahan,
		&i.TunjanganJabatan,
		&i.TunjanganLainPph21,
	)
	return i, err
}

const getNonStaffByNoPayroll = `-- name: GetNonStaffByNoPayroll :one
SELECT id, no_payroll, nama, foto, status_pph_21, alamat, nik, npwp, gender, jabatan, golongan, expat, ptkp, status_kerja, awal_kerja, akhir_kerja, gaji_pokok, tunjangan_kemahalan, tunjangan_perumahan, tunjangan_jabatan, tunjangan_lain_pph21 FROM non_staff WHERE no_payroll = ? LIMIT 1
`

func (q *Queries) GetNonStaffByNoPayroll(ctx context.Context, noPayroll string) (NonStaff, error) {
	row := q.db.QueryRowContext(ctx, getNonStaffByNoPayroll, noPayroll)
	var i NonStaff
	err := row.Scan(
		&i.ID,
		&i.NoPayroll,
		&i.Nama,
		&i.Foto,
		&i.StatusPph21,
		&i.Alamat,
		&i.Nik,
		&i.Npwp,
		&i.Gender,
		&i.Jabatan,
		&i.Golongan,
		&i.Expat,
		&i.Ptkp,
		&i.StatusKerja,
		&i.AwalKerja,
		&i.AkhirKerja,
		&i.GajiPokok,
		&i.TunjanganKemahalan,
		&i.TunjanganPerumahan,
		&i.TunjanganJabatan,
		&i.TunjanganLainPph21,
	)
	return i, err
}

const listNonStaff = `-- name: ListNonStaff :many
SELECT id, no_payroll, nama, foto, status_pph_21, alamat, nik, npwp, gender, jabatan, golongan, expat, ptkp, status_kerja, awal_kerja, akhir_kerja, gaji_pokok, tunjangan_kemahalan, tunjangan_perumahan, tunjangan_jabatan, tunjangan_lain_pph21 FROM non_staff
`

func (q *Queries) ListNonStaff(ctx context.Context) ([]NonStaff, error) {
	rows, err := q.db.QueryContext(ctx, listNonStaff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []NonStaff{}
	for rows.Next() {
		var i NonStaff
		if err := rows.Scan(
			&i.ID,
			&i.NoPayroll,
			&i.Nama,
			&i.Foto,
			&i.StatusPph21,
			&i.Alamat,
			&i.Nik,
			&i.Npwp,
			&i.Gender,
			&i.Jabatan,
			&i.Golongan,
			&i.Expat,
			&i.Ptkp,
			&i.StatusKerja,
			&i.AwalKerja,
			&i.AkhirKerja,
			&i.GajiPokok,
			&i.TunjanganKemahalan,
			&i.TunjanganPerumahan,
			&i.TunjanganJabatan,
			&i.TunjanganLainPph21,
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

func (q *Queries) CountNonStaff(ctx context.Context) int {
	var count int
	err := q.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM non_staff`).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (q *Queries) DeleteNonStaff(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, `DELETE FROM non_staff WHERE id = ?`, id)
	return err
}
