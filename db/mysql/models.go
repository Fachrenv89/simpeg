package db

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
}

type Lembur struct {
	ID        string  `json:"id"`
	NoPayroll string  `json:"no_payroll"`
	Nama      string  `json:"nama"`
	Jabatan   string  `json:"jabatan"`
	GajiPokok float32 `json:"gaji_pokok"`
	Basis     float32 `json:"basis"`
}

type LemburResponse struct {
	ID        string `json:"id"`
	NoPayroll string `json:"no_payroll"`
	Nama      string `json:"nama"`
	Jabatan   string `json:"jabatan"`
	GajiPokok string `json:"gaji_pokok"`
	Basis     string `json:"basis"`
}

type TanggalLembur struct {
	ID       string  `json:"id"`
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
	IDLembur string  `json:"id_lembur"`
}

type NonStaff struct {
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
